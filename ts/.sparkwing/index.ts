import { definePipeline, main } from "@sparkwing/sdk";
import { createHash } from "node:crypto";
const page = '<html><body><a href="/about.html">About</a></body></html>';
const about = '<html><body><a href="/index.html">Home</a></body></html>';
definePipeline({name:"HelloTs",plan(plan){
 const build=plan.job("build",async(ctx)=>{const digest=createHash("sha256").update(page).digest("hex");ctx.log.info(`build digest=${digest}`);return {digest,page};});
 plan.job("publish",async(ctx)=>{const out=await ctx.output<{digest:string,page:string}>(build);if(!out || out.digest!==createHash("sha256").update(out.page).digest("hex"))throw new Error("digest mismatch");ctx.log.info(`publish digest verified=${out.digest}`);}).needs(build);
}});
definePipeline({name:"Marketing",plan(plan){
 const build=plan.job("build",async(ctx)=>{ctx.log.info("marketing static site built pages=2");return {"/index.html":page,"/about.html":about};});
 const links=plan.job("link-check",async(ctx)=>{const pages=await ctx.output<Record<string,string>>(build);if(!pages)throw new Error("missing pages");let broken=0;for(const content of Object.values(pages)){for(const m of content.matchAll(/href="([^"]+)"/g)){if(!(m[1] in pages))broken++;}}ctx.log.info(`link-check broken links=${broken}`);if(broken)throw new Error("broken links");return {broken};}).needs(build);
 plan.job("preview",async(ctx)=>{const report=await ctx.output<{broken:number}>(links);if(!report || report.broken)throw new Error("invalid link report");ctx.log.info("marketing preview ready pages=2");}).needs(links);
}});
definePipeline({name:"ForkReader",secrets:[{name:"QA_SHARED",optional:true},{name:"QA_PIPELINE",optional:true}],plan(plan){plan.job("check-secrets",async(ctx)=>{for(const name of ["QA_SHARED","QA_PIPELINE"]){try{const value=await ctx.secret(name);ctx.log.info(`${name} present=${Boolean(value)} length=${value?.length??0}`);}catch(e){ctx.log.info(`${name} missing=${e}`);}}const jwt=await ctx.oidcToken("grid-fx");ctx.log.info(`oidc claims=${Buffer.from(jwt.split(".")[1],"base64url").toString()}`);});}});
definePipeline({name:"Fail",plan(plan){plan.job("fail",async()=>{throw new Error("grid fixture failure");});}});
await main();

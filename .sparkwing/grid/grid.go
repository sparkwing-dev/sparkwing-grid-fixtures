package grid
import (
 "context"
 "encoding/base64"
 "errors"
 "fmt"
 "strings"
 "time"
 "github.com/sparkwing-dev/sparkwing/sparkwing"
)
type Grid struct{sparkwing.Base}
type Sec struct{Shared string `sw:"QA_SHARED,optional"`; Pipeline string `sw:"QA_PIPELINE,optional"`}
func(Grid) Secrets() any{return &Sec{}}
func(Grid) Plan(_ context.Context,p *sparkwing.Plan,_ sparkwing.NoInputs,r sparkwing.RunContext)error{
 switch r.Pipeline{
 case "fx-parent":sparkwing.Job(p,"start-child",func(ctx context.Context)error{_,e:=sparkwing.RunAndAwait[any,sparkwing.NoInputs](ctx,"fx-child","greet",sparkwing.WithFreshBranch("grid-q-base"));return e}).StartsChildren("fx-child")
 case "fx-fail":sparkwing.Job(p,"fail",func(ctx context.Context)error{return errors.New("grid fixture failure")})
 case "fx-overlap":sparkwing.Job(p,"three-minute",func(ctx context.Context)error{sparkwing.Info(ctx,"interval start=%s",time.Now().UTC().Format(time.RFC3339Nano));select{case <-ctx.Done():return ctx.Err();case <-time.After(180*time.Second):};sparkwing.Info(ctx,"interval end=%s",time.Now().UTC().Format(time.RFC3339Nano));return nil}).Concurrency(sparkwing.NewConcurrencyGroup("fx-overlap",sparkwing.ConcurrencyLimit{Capacity:1,Scope:sparkwing.ScopeGlobal,OnLimit:sparkwing.Queue}))
 default:sparkwing.Job(p,"greet",func(ctx context.Context)error{
  for _,name:=range []string{"QA_SHARED","QA_PIPELINE"}{v,e:=sparkwing.Secret(ctx,name);if e!=nil{sparkwing.Info(ctx,"%s error=%v ErrSecretMissing=%v",name,e,errors.Is(e,sparkwing.ErrSecretMissing))}else{sparkwing.Info(ctx,"%s secret length=%d",name,len(v))}}
  t,e:=sparkwing.OIDCToken(ctx,"grid-fx");if e!=nil{return e};parts:=strings.Split(t,".");if len(parts)!=3{return fmt.Errorf("invalid token format")};b,e:=base64.RawURLEncoding.DecodeString(parts[1]);if e!=nil{return e};sparkwing.Info(ctx,"oidc claims=%s",b);sparkwing.Info(ctx,"hello grid fixture");return nil
 })
 };return nil
}
func init(){sparkwing.Register("Grid",func()sparkwing.Pipeline[sparkwing.NoInputs]{return Grid{}});for _,n:=range []string{"hello","fx-other","fx-pin","fx-parent","fx-child","fx-fail","fx-overlap","grid-q-actions","grid-q-embedded"}{sparkwing.Register(n,func()sparkwing.Pipeline[sparkwing.NoInputs]{return Grid{}})}}

import cpp
import coresInternal.OAI.Functions
import semmle.code.cpp.dataflow.DataFlow
import utils.Utils
class HttpMethodFlowConfig extends DataFlow::Configuration{
    HttpMethodFlowConfig(){ this = "HttpMethodFlowConfig"}

    override predicate isSource(DataFlow::Node node){
        exists(CoreServiceCall c, Variable v | v = c.getArgument(1).(FunctionCall).getArgument(_).(VariableAccess).getTarget()
                | node.asDefiningArgument() = v.getAnAccess())
    }

    override predicate isSink(DataFlow::Node node) {
        node.asExpr() = any(CoreServiceCall c).getArgument(1)
    }

    override predicate isBarrier(DataFlow::Node node){
        exists(FunctionCall f | f.getTarget().getName() = "operator=" and f.getQualifier().toString() = "method" 
            | node.asExpr() = f.getQualifier())
    }
}




// from CoreServiceCall s, Variable v, FunctionCall f
// where v = s.getArgument(1).(FunctionCall).getArgument(_).(VariableAccess).getTarget()
// and f = v.getAnAccess().getParent().(FunctionCall)
// and f.getTarget().getName() = ("operator=")
// and exists(ControlFlowNode c | c = f and c = s.getAPredecessor+())
// and not exists (ControlFlowNode c, FunctionCall f2 | c = f and 
//     f2 = s.getArgument(1).(FunctionCall).getArgument(_).(VariableAccess).getTarget().getAnAccess().getParent().(FunctionCall)
//     and f2.getTarget().getName() = "operator="
//     and f.getASuccessor+() = f2 and s.getAPredecessor+() = f2)
// // and f.getTarget().getName() = "= operator"
// // and src = DataFlow::exprNode(v.getAnAccess())
// // and sink = DataFlow::exprNode(s.getArgument(1))
// // and DataFlow::localFlow(src, sink)
// select f, s


// from HttpMethodFlowConfig cfg, DataFlow::Node src, DataFlow::Node sink
// where cfg.hasFlow(src, sink)
// select src, sink


// from DataFlow::Node src, DataFlow::Node sink
// where sink.asExpr() = any(CoreServiceCall c).getArgument(1).(FunctionCall).getAnArgument()
// and DataFlow::localFlowStep(src, sink)
// select unique(DataFlow::Node src2 | DataFlow::localFlowStep(src, sink) and src2 = src), sink

// from CoreServiceCall c
// where not exists(c.getHttpTypeCore())
// select c, c.getPathCore()

// from Function f, ControlFlowNode n
// where f.getName() = "ue_authentications_deregister_post"
// and n.getEnclosingElement() = f
// select f, n, n.(ExprStmt).getExpr().(FunctionCall).getArgument(_), n.(ExprStmt).getExpr().(FunctionCall).getArgument(_).getAQlClass()

// from EnumConstant e
// where e.getName() = "Not_Implemented"
// select e, e.getAnAccess()

from StringLiteral s, EnumConstantAccess e, CoreHandlers h, Function f, Function f2
where s.toString().matches("%not been implemented%")
and (e.getTarget().getName() = "Not_Implemented" or e.getTarget().getName() = "HTTP_STATUS_CODE_501_NOT_IMPLEMENTED")
// and (f.getName() = h.getName() or f.getACallToThisFunction().getEnclosingFunction().getName() = h.getName())
// and not (s.getEnclosingFunction() = f or e.getEnclosingFunction() = f)
and (e.getEnclosingFunction() = f or s.getEnclosingFunction() = f)
select f
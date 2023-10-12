private import BaseLanguage::BaseLanguageLib as BL
import Nodes

predicate hasLocalTaint(DataFlowNode src, DataFlowNode sink){
    BL::TaintTracking::localTaint(src, sink)
}
predicate hasLocalFlow(DataFlowNode src, DataFlowNode sink){
    BL::DataFlow::localFlow(src, sink)
}

DataFlowNode getExprNode(BL::Expr e){
    result = BL::DataFlow::exprNode(e)
}

DataFlowNode getParamNode(BL::Parameter p){
    result = BL::DataFlow::parameterNode(p)
}


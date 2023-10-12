private import BaseLanguage::BaseLanguageLib as BL
import Functions
private import Dataflow

class ErrorDef extends DataFlow::DataFlowNode {
    ErrorDef(){
        exists(CoreCallNode c | c.getAResult().getType() instanceof BL::ErrorType |
        this = c.getResult(_) and this.getType() instanceof BL::ErrorType)
        and
        exists(BL::DataFlow::PostUpdateNode p | DataFlow::hasLocalFlow(this, p))
        // this.getType() instanceof BL::ErrorType
        // and
        // exists(BL::Assignment a | a.getLhs(_) = this.asExpr())
        // and
        // getFunction().getName() = "GenerateAuthDataProcedure"
    }

    predicate isAssigned(){

        exists(BL::DataFlow::PostUpdateNode p | DataFlow::hasLocalFlow(this, p))
    }

    DataFlow::DataFlowNode getErrorNode(){
        exists(BL::DataFlow::PostUpdateNode p | DataFlow::hasLocalFlow(this, p) | result = p.getPreUpdateNode())
    }


}


//Abstract this out for the localTaint
//Need:
// localTaint wrapper
// Return definition (no error)
// dominates wrapper (for check -> return)
class CoreReturn extends BL::IR::ReturnInstruction {

    predicate hasReturnedExpr(){
        exists(BL::Expr e | this.getReturnStmt().getAnExpr() = e)
    }

    BL::Expr getAReturnedExpr(){
        result = this.getReturnStmt().getAnExpr()
    }

    string getLocationInfo(){
        result = getReturnStmt().getLocation().toString()
    }
}

class ErrorCheck extends BL::IfStmt {
    ErrorDef checkedError;

    ErrorCheck(){ 
        exists(ErrorDef ed, BL::Expr e | this.getCond().(BL::EqualityTestExpr) = e and
        ( ed.(BL::DataFlow::ExprNode).getExpr().getParent() =e
        or ed.getASuccessor().(BL::DataFlow::ExprNode).getExpr().getParent() = e
        or ed.getASuccessor().getASuccessor().(BL::DataFlow::ExprNode).getExpr().getParent() = e)
        | checkedError = ed)
    }

    BL::ControlFlow::Node getTrueGuardNode(){
        exists(BL::ControlFlow::ConditionGuardNode c | this.getCond() = c.getCondition()
            and c.ensures(DataFlow::getExprNode(this.getCond()), true) | result = c)
    }

    DataFlow::DataFlowNode getCondDFNode(){
        result = DataFlow::getExprNode(this.getCond())
    }

    string getLocationInfo(){
        result = getLocation().toString()
    }
}
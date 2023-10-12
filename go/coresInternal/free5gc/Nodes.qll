private import BaseLanguage::BaseLanguageLib as BL
import Functions
import AKA::AKA
import LoggerCall

class DataFlowNode extends BL::DataFlow::Node{

    predicate correspondsToExpr(BL::Expr e){
        this.asExpr() = e
    }

    predicate correspondsToParamOf(CoreFunction f) {
        this.asParameter() = f.getAParameter()
    }

    predicate correspondsToArgOf(CoreFunction f){
        this = f.getACall().getAnArgument()
    }

    string getParamFunction(){
        result = this.(BL::DataFlow::ParameterNode).asParameter().getFunction().getName().toString()
    }

    string getLocationInfo(){
        result = getFile().getRelativePath()+":"+getStartLine()+":"+getStartColumn()
    }

    predicate readsVar(AkaVariableBase v){ 
        this = v.getSourceVariable().getARead()
    }


    predicate writesVar(AkaVariableBase v){ 
        this = v.getSourceVariable().getAWrite().getRhs()
    }

    CoreFunction getFunction(){
        // exists(CoreFunction f, BL::Expr e | this.correspondsToExpr(e)
        //     and f.getBody() = e.getEnclosingFunction().getBody())

        result = this.getEnclosingCallable().asFunction()

        // result = this.asExpr().getEnclosingFunction().get
    }
}


class AkaVariableBase extends DataFlowNode, BL::DataFlow::SsaNode {

    AkaVariableBase(){
        not this.toString().matches("%phi(%")
    }
    override CoreFunction getFunction(){
        this.getEnclosingCallable().asFunction() = result
    }

    BL::Expr getDefExpr(){
        result = this.getSourceVariable().getDeclaration()
    }

}



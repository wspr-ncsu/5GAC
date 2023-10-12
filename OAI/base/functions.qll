import base.Target
import coreapi.Handler
import utils.Utils
class Function extends Target::CoreFunction {

    Call getCallInternal() {
        result = this.getCallCore()
    }

    Handler getHandler() {
        exists(Handler h | Utils::isDescendant(this, h.asFunction()) | result = h)
    }

    predicate canCauseCrash(){
        this.canCrashCore()
    }

    predicate recoversFromCrash(){
        this.recoversFromCrashCore()
    }

    predicate mayLeadToCrash(){
        exists(Function f | f.canCauseCrash() and Utils::isChild(this, f))
    }


}

class Call extends Target::CoreCall {
    Function getCallingFunction() {
        result = this.getCallerCore()
    }

    Function getCallTarget() {
        result = this.getTargetCore()
    }

}

class CallNode extends Target::CoreCallNode {

}

// class Variable extends Target::DataFlow::AkaVariableBase{
//     string getHashInfo(){
//         result = this.getLocationInfo()
//     }
// }



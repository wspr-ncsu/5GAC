private import BaseLanguage::BaseLanguageLib as BL
import Expr
// import Handlers
// import internalUtils
/**
 * Name: CoreFunction
 * Connects base language of target Core to cellular library.
 * Define a function belonging to the 5G Core in the marked areas.
 * 
 */


class CoreFunction extends BL::Function {


    /**
     * Core-specific logic to grab what NF this function belongs to
     */
    string getNFCore(){ result = this.getFile().getAbsolutePath().regexpCapture(".*/oai-cn5g-([A-Za-z0-9]+)/.*",1)}

    CoreCall getCallCore() {
        exists(CoreCall c | c.getTargetCore() = this | result = c)
    }

    CoreFunction getCallerCore() {
        result = getCallCore().getCallerCore()
    }

    predicate canCrashCore(){
        none()
    }

    predicate recoversFromCrashCore(){
        none()
    }

}

/**
 * Name: CoreCall
 * Connects base language function call to cellular library.
 * The CoreCall class is for a control flow, AST, or IR node.
 * CoreCallNode is the class to connect calls for dataflow.
 * 
 * The extension of this class (super class) will depend on the target language.
 * Make sure the extension is what you need. This class is primarily used internally to traverse call graphs.
 */

class CoreCall extends BL::Expr {
    CoreCall(){
        this instanceof BL::FunctionCall
        or
        this instanceof BL::FunctionAccess
    }
    CoreFunction getCallerCore(){
        result = this.getEnclosingFunction()
    }

    CoreFunction getTargetCore() {
        result = this.(BL::FunctionCall).getTarget()
        or
        result = this.(BL::FunctionAccess).getTarget()
    }

    BL::Expr getArgument(int n) {
        result = this.(BL::FunctionCall).getArgument(n)
    }

    int getNumArguments() {
        result = this.(BL::FunctionCall).getNumberOfArguments()
    }


}


class CoreCallNode extends BL::DataFlow::Node {


}

class CoreServiceCall extends CoreCall {
    
    CoreServiceCall() {
        exists(CoreFunction f | f.getName() = "curl_http_client" 
                                or f.getName() = "curl_create_handle" 
                                or f.getName() = "send_request"
                | f.getACallToThisFunction() = this)
    }

    HttpMethodArgument getHttpMethodArgument(){
        exists(HttpMethodArgument hma | hma.getServiceCall() = this | result = hma)
    }

    string getHttpTypeCore() {
        if exists(getHttpMethodArgument().getAssignedValue())
            then result = getHttpMethodArgument().getAssignedValue()
        else result = "UNKNOWN"
    }

    UriArgument getUriArgument(){
        exists(UriArgument ua | ua.getServiceCall() = this | result = ua)
    }
    string getPathCore(){
        result = getUriArgument().getUriString()
    }

    string getTargetNFCore(){
        result = "Post-Process From URI"
    }

    string getTargetHandlerCore(){
        result = "Unknown"
    }

    string getPathRootCore(){
        result = "Post-Process from URI"
    }

    CoreCall getCallCore(){
        result = this
    }

    string getSourceNFCore(){
        result = this.getCallerCore().getNFCore()
    }



}

class CoreTerminatingFunction extends CoreFunction{
    CoreTerminatingFunction(){
        this.getName() = ["main","ogs_sbi_discover_and_send"]
    }
}


class ApiRootFunctions extends CoreFunction{
    ApiRootFunctions(){
        this.getName().regexpMatch("get_[A-Za-z0-9]+_api_root")
    }

    string getApiRoot(){
        result = buildStringFromFunctionCall(this.getParameter(0).getAnAccess().getParent())
    }
}

class OpEqCall extends CoreCall{
    OpEqCall(){
        this.getTargetCore().getName() = "operator="
    }
}

class OpPlusCall extends CoreCall{
    OpPlusCall(){
        this.getTargetCore().getName() = "operator+"
    }
}

class OpPlusEqCall extends CoreCall{
    OpPlusEqCall(){
        this.getTargetCore().getName() = "operator+="
    }
}

class GetUriCall extends CoreCall{
    GetUriCall(){
        toString().regexpMatch("call to get_.+_uri")
    }
}
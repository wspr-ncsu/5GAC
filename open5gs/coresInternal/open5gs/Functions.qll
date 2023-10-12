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
    string getNFCore(){ result = "placeholder"}

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

class CoreServiceCall extends CoreFunction {
    

    CoreServiceCall() {
        exists(CoreFunction f | f.getName() = "ogs_sbi_build_request" | this.calls(f))
    }

    string getHttpTypeCore() {
        exists(SBIHttpAssignment sha | sha.getEnclosingFunction() = this | result = sha.getMethod())
    }

    string getPathRootCore(){
        exists(SBIVerAssignment ver, SBIServiceAssignment ser | 
            ver.getEnclosingFunction() = this and ser.getEnclosingFunction() = this |
            result = "{apiRoot}/"+ser.getService() + "/" + ver.getApiVersion())

    }

    private string buildPath(){
        result = concat(int i | i = [0 .. max(int num | num = any(SBIResourceAssignment res | res.getEnclosingFunction() = this | res.getArrayIndex()))] 
                | any(SBIResourceAssignment res | res.getEnclosingFunction() = this and res.getArrayIndex() = i | res.getResourceValue()), "/" order by i)
    }


    predicate hasPassedPath() {
        exists(BL::StringLiteral s , SBIServiceAssignment ser| 
        this.getCallCore().getParent().(CoreCall).getArgument(_) = s 
        and not s.getValue() = ser.getService() )
    }

    int getPassedPathIndex(){
        exists(BL::StringLiteral s , SBIServiceAssignment ser, int i | 
        i = [0 .. this.getCallCore().getParent().(CoreCall).getNumArguments()]
        and this.getCallCore().getParent().(CoreCall).getArgument(i) = s 
        and not s.getValue() = ser.getService() |
        result = i)
    }

    string getPassedPath(){
        exists(BL::StringLiteral s , SBIServiceAssignment ser| 
            this.getCallCore().getParent().(CoreCall).getArgument(_) = s 
            and not s.getValue() = ser.getService() |
            result = s.getValue())
    }

    string getPassedPathParameterName(){
        result = this.getCallCore().getParent().(CoreCall).getTargetCore().getParameter(getPassedPathIndex()).getName()
    }


    string getPathCore(){
        if hasPassedPath() then (
            result = buildPath().replaceAll(this.getPassedPathParameterName(), getPassedPath())
            // result = "aaaaaaaaaa"
        )
        // else result = "bbbbbbbbbb"
        else (
            result = buildPath()
            // result = "bbbbbbbb"
        )
        
    }



    string getTargetNFCore() {
        result = this.getName().regexpCapture("([a-zA-Z0-9]+)_n([a-zA-Z0-9]+).*", 2)
    }

    string getSourceNFCore(){
        result = this.getName().regexpCapture("([a-zA-Z0-9]+)_n([a-zA-Z0-9]+).*", 1)
    }

    CoreHandlers getTargetHandlerCore() {
        result = any(CoreHandlers h)
    }


    // bindingset[path]
    string getColonizedPathCore(){
        result = getPathCore()
    }

    string getPlmnIdPathCore(){
        result = getPathCore()
    }


    // CoreRouter getServiceCallRouterCore() {
    // }

    BL::FunctionCall getServiceCallInternal(){
        result = this.getACallToThisFunction()
    }

    BL::FunctionAccess getServiceAccessInternal(){
        result = this.getAnAccess()
    }

    override CoreCall getCallCore(){
        result = getServiceAccessInternal() 
        or
        result = getServiceCallInternal()
    }

    predicate isCalled(){
        (exists(BL::FunctionAccess a | a = this.getAnAccess()) or exists(BL::FunctionCall c | c = this.getACallToThisFunction()))
    }
}

class CoreTerminatingFunction extends CoreFunction{
    CoreTerminatingFunction(){
        this.getName() = ["main","ogs_sbi_discover_and_send"]
    }
}
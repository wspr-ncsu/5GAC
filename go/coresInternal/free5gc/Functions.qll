private import BaseLanguage::BaseLanguageLib as BL
import Handlers
import internalUtils

/**
 * Name: CoreFunction
 * Connects base language of target Core to cellular library.
 * Define a function belonging to the 5G Core in the marked areas.
 * 
 */
class CoreFunction extends BL::Function {

    /** Add Definition of core function here. */
    CoreFunction() { this instanceof BL::Function}


    /**
     * Core-specific logic to grab what NF this function belongs to
     */
    string getNFCore(){ result = "a"}

    CoreCall getCallCore() {
        result = any(CoreCall c | c.getTarget() = this)
    }

    predicate isInThirdPartyLibrary(){
        if removePackageBase(getPackage().toString()).regexpMatch(".+\\..+") then 
            not removePackageBase(getPackage().toString()).regexpMatch(".*free5gc.*")
        else
            not this.getPackage().toString().regexpMatch(".*free5gc.*")  
    }

    predicate canCrashCore(){
        exists(CoreFunction f | f.mustPanic() | f.getCallCore().getEnclosingFunction() = this.getFuncDecl())
    }

    predicate recoversFromCrashCore(){
        exists(BL::DeferStmt d | d.getEnclosingFunction+() = this.getFuncDecl() 
                | (isChild(this, BL::Builtin::recover()) 
                    or BL::Builtin::recover().getAReference().getEnclosingFunction+() = this.getFuncDecl()))
    }

    string getHashInfo(){
        result = this.getQualifiedName()
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

class CoreCall extends BL::CallExpr {
    CoreFunction getCallerCore(){
        exists(CoreFunction f | this.getEnclosingFunction().getBody() = f.getBody() | result = f  )
    }

    CoreFunction getTargetCore() {
        result = this.getTarget()
    }

    string getHashInfo(){
        result = this.getLocation().toString()
    }

}


class CoreCallNode extends BL::DataFlow::CallNode{

    string getHashInfo(){
        result = this.getEndLine()+this.getRoot().getLocation().toString()
    }

}

class CoreServiceCall extends CoreFunction {
    BL::CallExpr apicall;

    CoreServiceCall() {
        exists(BL::Function f, BL::CallExpr c | 
            c.getCalleeName() = "CallAPI" and 
            c.getEnclosingFunction().getBody() = f.getBody() |
            this = f and apicall = c )
    }


    predicate isCalled(){
        exists(BL::DataFlow::CallNode c | c.getTarget() = this)
    }

    string getHttpTypeCore() {
        exists(
            BL::ValueSpec v | 
            v.getAName().toLowerCase() = "localVarHTTPMethod".toLowerCase() and 
            v.getEnclosingFunction() = apicall.getEnclosingFunction() | 
            result = v.getAnInit().(BL::CallExpr).getAnArgument().getExactValue()
        )
    }

    string getPathRootCore(){

        exists( CoreFunction f, BL::KeyValueExpr k, BL::DefineStmt d | 
            (f.hasQualifiedName(this.getPackage().getPath(), "NewConfiguration") and 
            k.getEnclosingFunction().getBody() = f.getBody() and
            k.getKey().toString() = "url") 
            and
            (d.getLhs().toString() = "localVarPath" and
            d.getEnclosingFunction() = apicall.getEnclosingFunction())|
             result = (k.getValue().getExactValue())
        )
        // and
        // this.getName() = "GetNFInstances"
    }




    string getPathCore(){
        exists(BL::DefineStmt d | 
            d.getEnclosingFunction() = apicall.getEnclosingFunction() and
            d.getLhs().toString() = "localVarPath" | 
            result = (d.getRhs().(BL::AddExpr).getRightOperand().getExactValue()))
    }


    string getTargetNFCore() {
        exists(CoreRouter r, CoreHandlers h | 
            (r.getPathCore() = this.getPathCore() 
            or r.getPathCore() = this.getPlmnIdPathCore())
            and
            h.getName() = r.getHandlerNameCore() 
            and r.getRouterRootPath() = this.getPathRootCore().replaceAll("{apiRoot}", "") |
            result = matchNF(r.getFile().getRelativePath().splitAt("/")))
        or
        exists(CoreRoutingFunction r, CoreHandlers h | 
            (r.getPathCore() = this.getPathCore() 
            or r.getPathCore() = this.getPlmnIdPathCore())
            and
            h.getName() = r.getHandlerNameCore() 
            and r.getRouterRootPath() = this.getPathRootCore().replaceAll("{apiRoot}", "") |
            result = matchNF(r.getFile().getRelativePath().splitAt("/")))
    }

    CoreHandlers getTargetHandlerCore() {
        exists(CoreRouter r, CoreHandlers h | 
            (r.getPathCore() = this.getPathCore() or r.getPathCore() = this.getPlmnIdPathCore()) and
            h.getName() = r.getHandlerNameCore() 
            and r.getHttpTypeCore() = this.getHttpTypeCore() 
            and r.getRouterRootPath() = this.getPathRootCore().replaceAll("{apiRoot}", "")
            | result = h)
        or
        exists(CoreRoutingFunction r, CoreHandlers h | 
            (r.getPathCore() = this.getPathCore() or r.getPathCore() = this.getPlmnIdPathCore()) and
            h.getName() = r.getHandlerNameCore() 
            and r.getHttpTypeCore() = this.getHttpTypeCore() 
            and r.getRouterRootPath() = this.getPathRootCore().replaceAll("{apiRoot}", "")
            | result = h)
    }


    // bindingset[path]
    string getColonizedPathCore(){
        result = this.getPathCore()
        .regexpReplaceAll("\\{(.[^\\}]*)\\}", ":$1")
    }

    string getPlmnIdPathCore(){
        result = this.getColonizedPathCore()
        .replaceAll("authentication-data", ":servingPlmnId") // Handle annoying path change in routers
        .replaceAll("context-data", ":servingPlmnId")  // Handle annoying path change in routers 
    }


    CoreRouter getServiceCallRouterCore() {
        result = any(CoreRouter r 
            | r.getRouterRootPath() = this.getPathRootCore().replaceAll("{apiRoot}", "") 
            and r.getPathCore() = this.getPlmnIdPathCore() | r)

    }


    string getSourceNFCore(){
        result = matchNF(getBody().getFile().getRelativePath().splitAt("/"))
    }
}

class CoreTerminatingFunction extends CoreFunction {
    CoreTerminatingFunction() {
        hasQualifiedName(_, ["Start", "Terminate"])
        
    }
}
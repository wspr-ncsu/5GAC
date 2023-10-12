import Functions
import Routers
import internalUtils
import generatedLibs

private class HandlerDef extends GeneratedHandlers{ 
}


library class CoreHandlers extends CoreFunction {

    CoreHandlers(){
        this instanceof CoreHandlersRouteFunction
        or
        this instanceof CoreHandlersRouter
        or
        this instanceof CoreGMMHandler
    }

    
    string getHttpTypeInternal() {
        // exists(CoreRouter r | r.getHandlerNameCore() = this.getName() | 
        // result = r.getHttpTypeCore())
        result = this.(CoreHandlersRouteFunction).getHttpTypeCore()
        or
        result = this.(CoreHandlersRouter).getHttpTypeCore()
        or
        result = this.(CoreGMMHandler).getHttpTypeCore()
        // result = "debug"
    }


    /**
     * Core-specific logic to determine what URL reaches this handler.
     * Returns: String of handler URL
     */
    string getPathInternal(){         
        // exists(CoreRouter r | r.getHandlerNameCore() = this.getName() | 
        // result = r.getPathCore())
        result = this.(CoreHandlersRouteFunction).getPathCore()
        or
        result = this.(CoreHandlersRouter).getPathCore()
        or
        result = this.(CoreGMMHandler).getPathCore()
    }

    string getNFInternal() {result = this.getNFCore() or result = this.(CoreGMMHandler).getNFCore()}

    CoreHandlers asFunction(){
        result = this
    }
}

library class CoreHandlersRouter extends CoreFunction{
    // string url;
    // string nf;
    // string name;
    // string http_type;
    CoreRouter router;

    CoreHandlersRouter() {
        exists(CoreRouter r | r.getHandlerNameCore() = this.getName() | router = r)
    }

    /**
     * Core-specific logic to determine what HTTP type reaches this handler
     * Returns: String of Handler HTTP type
     */
    string getHttpTypeCore() {
        // exists(CoreRouter r | r.getHandlerNameCore() = this.getName() | 
        // result = r.getHttpTypeCore())
        result = router.getHttpTypeCore()
        // result = "debug"
    }


    /**
     * Core-specific logic to determine what URL reaches this handler.
     * Returns: String of handler URL
     */
    string getPathCore(){         
        // exists(CoreRouter r | r.getHandlerNameCore() = this.getName() | 
        // result = r.getPathCore())
        result = router.getPathCore()
    }

    // string getPathRootCore(){
    //     exists(CoreServiceCall sc | sc.getTargetHandlerCore() = this | result = sc.getPathRootCore())
    // }


    /**
     * Core-specific logic to determine what NF this handler belongs to.
     * Overrides getNFInternal() from CoreFunction. Only define if handler is different/required
     * Returns: String of handler NF
     */
    override string getNFCore() { result = matchNF(router.getFile().getRelativePath().splitAt("/"))}

}

library class CoreHandlersRouteFunction extends CoreFunction{
    // string url;
    // string nf;
    // string name;
    // string http_type;
    CoreRoutingFunction routing;

    CoreHandlersRouteFunction() {
        exists(CoreRoutingFunction r | r.getHandlerNameCore() = this.getName() | routing = r)
    }

    /**
     * Core-specific logic to determine what HTTP type reaches this handler
     * Returns: String of Handler HTTP type
     */
    string getHttpTypeCore() {
        // exists(CoreRouter r | r.getHandlerNameCore() = this.getName() | 
        // result = r.getHttpTypeCore())
        result = routing.getHttpTypeCore()
        // result = "debug"
    }


    /**
     * Core-specific logic to determine what URL reaches this handler.
     * Returns: String of handler URL
     */
    string getPathCore(){         
        // exists(CoreRouter r | r.getHandlerNameCore() = this.getName() | 
        // result = r.getPathCore())
        result = routing.getPathCore()
    }

    // string getPathRootCore(){
    //     exists(CoreServiceCall sc | sc.getTargetHandlerCore() = this | result = sc.getPathRootCore())
    // }


    /**
     * Core-specific logic to determine what NF this handler belongs to.
     * Overrides getNFInternal() from CoreFunction. Only define if handler is different/required
     * Returns: String of handler NF
     */
    override string getNFCore() { result = matchNF(routing.getFile().getRelativePath().splitAt("/"))}

}

library class CoreGMMHandler extends CoreFunction {

    // Match names for now, can find from the state machine definition/intialization, or ideally somewhere in specs
    CoreGMMHandler(){
        this.getName() = ["DeRegistered", "Authentication", "SecurityMode", "ContextSetup", "Registered", "DeregisteredInitiated"]
    }

    override string getNFCore(){result = matchNF(getBody().getFile().getRelativePath().splitAt("/"))}

    string getPathCore(){result = "Message from UE"}

    string getHttpTypeCore(){result = "N/A"}
}
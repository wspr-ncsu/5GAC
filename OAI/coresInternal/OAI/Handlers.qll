import Functions
import Expr
private import BaseLanguage::BaseLanguageLib as BL

class Router extends BL::Function {
    Router(){
        this.getNamespace().getName() = "Routes" and not this.getName() = "bind"
    }
}

class RouteSetup extends BL::FunctionCall{
    RouteSetup(){
        this.getTarget().getNamespace().getName() = "Routes"
        and
        not this.getTarget().getName() = "bind"
    }

    string getRouteBase(){
        exists(BL::FunctionCall basic_string | 
            basic_string.getTarget().getName() = "basic_string" 
            and
            getAnArgument+().(BL::VariableAccess).getTarget().getAnAssignedValue() = basic_string |
            result = basic_string.getAnArgument().toString()
            )
    }

    string getRouteTarget(){
        result = concat(getAnArgument+().(BL::StringLiteral).getValue())
    }

    string getRouteVersion(){
        if exists(BL::ValueFieldAccess v | getAnArgument+() = v | v.toString().matches("%version%"))
        then result = "v1"
        else result = ""
    }

    string getRoute(){
        result = getRouteBase() + getRouteVersion() + getRouteTarget()
    }

    CoreCall getHandler(){
        result = getArgument(2).getAChild().(CoreCall)
    }

    string getHttpType(){
        result = this.getTarget().getName()
    }
}



class CoreHandlers extends CoreCall {

    RouteSetup route;


    CoreHandlers(){
        exists(RouteSetup r | this = r.getHandler() and route = r)

    }

    
    string getHttpTypeInternal() {
        result = route.getHttpType()
    }


    string getSerivceName(){
        result = route.getRouteBase().regexpCapture("(/?[0-9A-Za-z-]*)/?.*", 1)
    }

    CoreFunction asFunction(){
        result = this.getTargetCore()
    }

    // /**
    //  * Core-specific logic to determine what URL reaches this handler.
    //  * Returns: String of handler URL
    //  */
    string getPathInternal(){
        result = route.getRoute()
    }

    string getNFInternal(){result = this.getTargetCore().getNFCore()}

    string getName(){
        result =this.getTargetCore().getName()
    }


    
}
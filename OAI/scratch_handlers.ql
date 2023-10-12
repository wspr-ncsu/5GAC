import cpp

class Router extends Function {
    Router(){
        getNamespace().getName() = "Routes" and not getName() = "bind"
    }
}

class RouteSetup extends FunctionCall{
    RouteSetup(){
        getTarget() instanceof Router
    }

    string getRouteBase(){
        exists(FunctionCall basic_string | 
            basic_string.getTarget().getName() = "basic_string" 
            and
            getAnArgument+().(VariableAccess).getTarget().getAnAssignedValue() = basic_string |
            result = basic_string.getAnArgument().toString()
            )
    }

    string getRouteTarget(){
        result = concat(getAnArgument+().(StringLiteral).getValue())
    }

    string getRouteVersion(){
        if exists(ValueFieldAccess v | getAnArgument+() = v | v.toString().matches("%version%"))
        then result = "v1"
        else result = ""
    }

    string getRoute(){
        result = getRouteBase() + getRouteVersion() + getRouteTarget()
    }

    Function getHandler(){
        result = getArgument(2).getAChild().(FunctionAccess).getTarget()
    }
}



// from Function f, FunctionCall c, Expr e
// where f.getNamespace().getName() = "Routes" and c.getTarget() = f and not f.getName() = "bind"
// and e = c.getArgument(1).getAChild()
// select c, e, e.(VariableAccess).getTarget().getAnAssignedValue().(FunctionCall).getAnArgument()



from Function f, FunctionCall c, Expr e
where f.getNamespace().getName() = "Routes" and c.getTarget() = f and not f.getName() = "bind"
and e = c.getArgument(2)
select c, e.getAChild().(FunctionAccess).getTarget()

// from RouteSetup r
// where not exists(r.getRouteTarget())
// select r
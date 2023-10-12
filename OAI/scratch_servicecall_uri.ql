import cpp
import Handlers
import Functions
import Expr


Expr walkToEarliestAssignment(Expr e){
    if e instanceof VariableAccess
        then result = e.(VariableAccess).getTarget().getInitializer().getExpr().getAChild()
    else if e instanceof FunctionCall
        then result = e.(FunctionCall).getArgument(_)
    else result = e
}

FunctionCall getNextArg(Expr e){
    result = e.(FunctionCall).getArgument(0)
}

int getLineColNum(Location l){
    result = (l.getStartLine().toString() + l.getStartColumn().toString()).toInt()
}

string buildStringFromFunctionCall(Expr f){
    result = concat(Expr child | child = walkArgs*(f) and 
                (child instanceof ReferenceFieldAccess 
                or child instanceof StringLiteral 
                or child instanceof ValueFieldAccess
                or child instanceof PointerFieldAccess
                or child instanceof VariableAccess
                and not child instanceof Qualifier) 
            | child.toString() order by getLineColNum(child.getLocation()))
}

/// Parameter Arg Targets are events, don't care

/// getValue() gives operator+= strings

Expr walkArgs(FunctionCall e){
    result = e.getArgument(_)
}

// from UriArgument arg, OpPlusEqCall op
// where not exists(arg.getInitializedString())
// and not exists(arg.getAssignString())
// and op.(FunctionCall).getQualifier().(VariableAccess).getTarget().getAnAccess() = arg
// select op, arg
// // from ApiRootFunctions f
// select f, f.getParameter(0).getAnAccess()

from UriArgument arg
// where not exists()
select arg, arg.getUriString()
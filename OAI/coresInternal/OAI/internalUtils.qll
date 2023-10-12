private import cpp
private import Functions

bindingset [s]
int getUnspecCount(string s){
    result = count(string str 
            | str = s.splitAt("/")
                and str = "unspecified"
            | str )
}

Expr walkArgs(FunctionCall e){
    result = e.getArgument(_)
}

int getLineColNum(Location l){
    result = (l.getStartLine().toString() + l.getStartColumn().toString()).toInt()
}

language[monotonicAggregates]
string buildStringFromFunctionCall(Expr f){
    result = strictconcat(Expr child | child = walkArgs*(f) and 
                (child instanceof ReferenceFieldAccess 
                or child instanceof StringLiteral 
                or child instanceof ValueFieldAccess
                or child instanceof PointerFieldAccess
                or child instanceof VariableAccess
                and not child instanceof Qualifier)
                // and exists(string str | str = parseExpr(child) | s = str)
            | parseExpr(child) order by getLineColNum(child.getLocation()))
}

private string parseExpr(Expr child){
    // if child instanceof VariableAccess
    //     then if getOtherAccess(child).getParent() instanceof FunctionCall
    //     then exists(ApiRootFunctions f | getOtherAccess(child).getParent() = f.getACallToThisFunction() | result = f.getApiRoot())
    //     else result = child.toString()
    // else if child instanceof StringLiteral 
    //     then result = child.toString()
    // else result = "END"
    if child instanceof StringLiteral
        then result = child.toString()
    else if exists(getOtherAccess(child).getParent().(FunctionCall).getTarget().(ApiRootFunctions))
        then result = getOtherAccess(child).getParent().(FunctionCall).getTarget().(ApiRootFunctions).getApiRoot()
    // else if exists(walkArgs+(getOtherAccess(child).getParent().(OpEqCall)).(StringLiteral))   
    //     then result = buildStringFromFunctionCall(getOtherAccess(child).getParent().(OpEqCall).getArgument(0))
    else if exists(walkArgs+(child.(VariableAccess).getTarget().getInitializer().getExpr().(FunctionCall)).(StringLiteral))
        then result = buildStringFromFunctionCall(child.(VariableAccess).getTarget().getInitializer().getExpr().(FunctionCall))
    else result = "$VAR__" + child.toString()


}


VariableAccess getOtherAccess(VariableAccess var){
    result = var.getTarget().getAnAccess()
}

// Expr getOpEqAssign(VariableAccess var){
//     result = 
// }
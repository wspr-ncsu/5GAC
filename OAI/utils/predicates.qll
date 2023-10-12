import base.Core

private class NFNames extends string {
    NFNames() {
        this = "amf" or
        this = "ausf" or
        this = "n3iwf" or
        this = "nrf" or
        this = "nef" or
        this = "nssf" or
        this = "pcf" or
        this = "smf" or
        this = "udm" or
        this = "udr" or
        this = "upf"
    }

}


Core::Function getParent(Core::Function f){
    result = f.getCallInternal().getCallingFunction()

}

Core::Function getChildren(Core::Function f){
    result = any(Core::Call c | c.getCallingFunction() = f | c.getCallTarget())
}

// bindingset[s]
// predicate isDescendant(Core::Function child, string s){
//     getParent*(child).getName().matches(s)
// }

predicate isDescendant(Core::Function child, Core::Function parent){
    getParent*(child) = parent
}


bindingset[s]
string matchNF(string s){
     result = any(NFNames n | s.matches(n.toString()))
}

predicate isChild(Core::Function parent, Core::Function child){
    getChildren*(parent) = child
}
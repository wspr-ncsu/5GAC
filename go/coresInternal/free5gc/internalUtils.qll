import Functions

private class NFNames extends string {
    NFNames() {
        this = "amf" or
        this = "ausf" or
        this = "n3iwf" or
        this = "nrf" or
        this = "nssf" or
        this = "pcf" or
        this = "smf" or
        this = "udm" or
        this = "udr" or
        this = "upf"
    }

}

bindingset[s]
string matchNF(string s){
     result = any(NFNames n | s.matches(n.toString()) and result = n)
}

bindingset[p1,p2]
predicate comparePaths(string p1, string p2){
    p1.replaceAll("\"", "") = p2.replaceAll("\"","")
}

predicate isChild(CoreFunction parent, CoreFunction child){
    getChildren*(parent) = child
}

CoreFunction getChildren(CoreFunction f){
    result = any(CoreCall c | c.getCallerCore() = f | c.getTargetCore())
}

bindingset[pack]
string removePackageBase(string pack){
    result = concat(string ss, int i | ss = pack and i = [3 .. count(string s | s = pack.splitAt("/") | s)] | ss.splitAt("/", i), "/" )
}
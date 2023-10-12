import base.Target


class CFNode extends Target::ControlFlow::Node{
}


predicate checkDomination(CFNode dominator, CFNode dominee){
    Target::ControlFlow::checkDomination(dominator, dominee)
}

class TerminatingFunction extends Target::CoreTerminatingFunction{
    // TerminatingFunction(){
    // }
}

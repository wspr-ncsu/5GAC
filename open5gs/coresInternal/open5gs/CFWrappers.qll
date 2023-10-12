private import BaseLanguage::BaseLanguageLib as BL

predicate checkDomination(Node dominator, Node dominee){
    BL::dominates(dominator, dominee)
}

class Node extends BL::ControlFlowNode{}
private import BaseLanguage::BaseLanguageLib as BL

predicate checkDomination(Node dominator, Node dominee){
    dominator.dominatesNode(dominee)
}

class Node extends BL::ControlFlow::Node{}
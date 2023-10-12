// import base.Core
// import coreapi.CoreAPI
// import base.Target
// import utils.Utils

// class AkaHandler extends CoreAPI::Handler {
//     AkaHandler(){
//         this = any(Target::AKA::GeneratedAkaHandlers a)
//     }
// }


// class AkaFunction extends Core::Function {
//     AkaHandler handler;
//     Core::Function parent;
//     string infotest;

//     AkaFunction(){
//         exists(AkaHandler h | Utils::isDescendant(this, h) | handler = h)
//         and
//         not isInThirdPartyLibrary()
//         and 
//         parent = Utils::getParent(this)
//         and
//         infotest = this.getHashInfo()
//     }
// }

// class AkaVariable extends Core::Variable {

//     AkaVariable(){
//         exists(AkaFunction a | this.getFunction() = a)
//     }
// }

// class AuthVariable extends AkaVariable {

//     AuthVariable(){
//         exists(Core::AuthVarDiscovery discover | discover.hasFlowTo(this))

//     }

// }
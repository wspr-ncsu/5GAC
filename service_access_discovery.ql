import coreapi.CoreAPI
import utils.Utils
import base.Core

from CoreAPI::ServiceCall s, string handleName
// where s.getHandler() = h abnd
// where handleName = "debug"
where (exists(CoreAPI::Handler hh | Utils::isDescendant(s, hh.asFunction()) | handleName = hh.getName())
or (exists(Core::TerminatingFunction f | Utils::isDescendant(s, f) | handleName = f.getName()) 
        and not exists(CoreAPI::Handler hh | Utils::isDescendant(s, hh.asFunction()))
        ))
select s.getHttpType() as http_type,
handleName as serviceCallCFSource, 
s.getSourceNF() as sourceNF,
// h.getPath() as callerHandlerPath,
s.getName() as serviceCallFunction,
s.getTargetNF() as targetNF,
// s.getTargetHandler() as TargetHandler,
s.getColonizedPath() as targetPath,
s.getPathRoot() as targetRootPath
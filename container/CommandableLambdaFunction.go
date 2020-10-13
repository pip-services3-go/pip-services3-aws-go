package container

// /*@module container */
// import { ICommandable } from 'pip-services3-commons-node';
// import { CommandSet } from 'pip-services3-commons-node';
// import { Parameters } from 'pip-services3-commons-node';

// import { LambdaFunction } from './LambdaFunction';

// /*
//  Abstract AWS Lambda function, that acts as a container to instantiate and run components
//  and expose them via external entry point. All actions are automatically generated for commands
//  defined in [[https://rawgit.com/pip-services-node/pip-services3-commons-node/master/doc/api/interfaces/commands.icommandable.html ICommandable components]]. Each command is exposed as an action defined by "cmd" parameter.
//   
//  Container configuration for this Lambda function is stored in "./config/config.yml" file.
//  But this path can be overriden by CONFIG_PATH environment variable.
//  
//  ### Configuration parameters ###
//  
//  - dependencies:
//      - controller:                  override for Controller dependency
//  - connections:                   
//      - discovery_key:               (optional) a key to retrieve the connection from [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/connect.idiscovery.html IDiscovery]]
//      - region:                      (optional) AWS region
//  - credentials:    
//      - store_key:                   (optional) a key to retrieve the credentials from [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/auth.icredentialstore.html ICredentialStore]]
//      - access_id:                   AWS access/client id
//      - access_key:                  AWS access/client id
//  
//  ### References ###
//  
//  - \*:logger:\*:\*:1.0            (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/log.ilogger.html ILogger]] components to pass log messages
//  - \*:counters:\*:\*:1.0          (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/count.icounters.html ICounters]] components to pass collected measurements
//  - \*:discovery:\*:\*:1.0         (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/connect.idiscovery.html IDiscovery]] services to resolve connection
//  - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials
//  
//  See [[LambdaClient]]
//  
//  ### Example ###
//  
//      class MyLambdaFunction extends CommandableLambdaFunction {
//          private _controller: IMyController;
//          ...
//          public constructor() {
//              base("mygroup", "MyGroup lambda function");
//              this._dependencyResolver.put(
//                  "controller",
//                  new Descriptor("mygroup","controller","*","*","1.0")
//              );
//          }
//      }
//  
//      let lambda = new MyLambdaFunction();
//      
//      service.run((err) => {
//          console.log("MyLambdaFunction is started");
//      });
//  */
// export abstract class CommandableLambdaFunction extends LambdaFunction {

//     /*
//      Creates a new instance of this lambda function.
//      
//      - name          (optional) a container name (accessible via ContextInfo)
//      - description   (optional) a container description (accessible via ContextInfo)
//      */
//     public constructor(name: string, description?: string) {
//         super(name, description);
//         this._dependencyResolver.put('controller', 'none');
//     }

//     private registerCommandSet(commandSet: CommandSet) {
//         let commands = commandSet.getCommands();
//         for (let index = 0; index < commands.length; index++) {
//             let command = commands[index];

//             this.registerAction(command.getName(), null, (params, callback) => {
//                 let correlationId = params.correlation_id;
//                 let args = Parameters.fromValue(params);
//                 let timing = this.instrument(correlationId, this._info.name + '.' + command.getName());

//                 command.execute(correlationId, args, (err, result) => {
//                     timing.endTiming();
//                     callback(err, result);
//                 })
//             });
//         }
//     }

//     /*
//      Registers all actions in this lambda function.
//      */
//     public register(): void {
//         let controller: ICommandable = this._dependencyResolver.getOneRequired<ICommandable>('controller');
//         let commandSet = controller.getCommandSet();
//         this.registerCommandSet(commandSet);
//     }
// }

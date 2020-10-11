# timebox

Timebox is an Event Sourcing reference implementation, written in Go. Included is a simple Event Store and APIs for composing very basic Event Sourced services. It can either serve as the foundation for something bigger, or serve as an example of some key concepts. It is not, however, a stand-alone framework.

## Event Sourcing: A Glossary

* Event Sourcing: The practice of describing all changes to a system as a sequence of Events rather than as operations against mutable data models.

* Aggregate: A group of objects or entities that are related to one another. They are accessible via a single object or entity, known as the Aggregate Root. An event sourced system performs all Commands through an Aggregate, rather than directly against its descendants. The Aggregate establishes and maintains the integrity of its descendants.

* Applier (Calculator): An interface that is implemented in order to Apply a raised Event to an Aggregate. Appliers only apply state changes to an Aggregate. They do not verify the validity of the Aggregate or the Event.

* Command: A message that encapsulates the intent of a user or component to perform some business-related action against the system. Command should be described in terms of the intent rather than the resulting change (ex: `CancelRide` rather than `SetRideStatus`). Depending on the current state of the system, a Command may or may not succeed.

* Command Handler: An interface that is implemented in order to handle a Command.

* Command Retry Handler: Wraps another Command handler, repeating the entire Command should a Version Consistency Error occur.

* Events: A message that describes a discrete change to a system. Events should generally be described in terms of business processes rather than implementation details (ex: `PassengerEnteredVehicle` rather than `DriverPushedButton`). Events must be applied in the order they were persisted.

* Event Handler: An interface that is implemented in order to handle an Event.

* Event Store: Stores the sequential set of Events that would be used to Hydrate an Aggregate, and the Raised Events created when successfully performing a Command against an Aggregate.

* Hydrate: To Apply the set of Events necessary to bring an Aggregate to a desired state. Most of the time, that state will be *current*, and so all stored Events for will be Applied.

* Typed Message Handler: A Message Handler designed to dispatch the handling of a Message to a specific Handler based on the Message's type. Multiple handlers may be registered per type. The dispatcher invokes the most recently registered handler first. Intermediate errors short-circuit the processing.

* Typed Command Handler: A Command Handler designed to relay Command handling to a specific Handler based on the Command's type. Only a single handler may be registered per type.

## Event Sourcing: One Possible Flow

1. Typed Command Handler accepts a Command and invokes the appropriate Command Handler for the Command's type.

2. Command Handler retrieves the Events necessary to hydrate an object graph based on the Aggregate ID provided by the Command (or a newly generated ID). It instantiates a skeletal graph and then begins to apply the Events to it.

3. Command Handler performs a check on the state of the hydrated graph in order to determine whether it can process the Command. If it can't, it returns a Failure. **Go To 7**.

4. The Command Handler raises New Events in order to reproduce the intended final state of the object graph.

5. (Optional) The Command Handler applies those Events to the graph, and performs a final integrity check. If the integrity check fails, a non-recoverable exception occurs. This would be considered a programmer error. **Go To 7**

6. Command Handler attempts to write the newly generated Events to the Store.

   6a. If the Store sees that the underlying Event stream shows a different Highest Version than previously noted, it will raise a recoverable Version Inconsistency Exception. That exception should be handled by an upstream Command Retry Handler, which will attempt to re-try the Command Handler. **Go To 2**

   6b. The Store persists the new Events, and the Command Handler returns Success.

7. Typed Command Handling Completed. Return Success or Failure.

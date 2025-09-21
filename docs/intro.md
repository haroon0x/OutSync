  ### 1. The Project: OutSync (What, Why, and How)


  What is it?
  OutSync is a small but powerful backend service whose only job is to guarantee that an "event" is sent to a message queue (like Kafka) whenever a critical action happens in our database.


  Why do we need it? (The Problem)
  Imagine a simple application where a user signs up. The app needs to do two things:
  1.  Save the new user to the database.
  2.  Send a "Welcome" email.


  The naive way is to do them in order. But what if your application crashes after saving the user but before sending the email? You now have a user who is registered but never got their welcome email. The
  system is in an inconsistent state. This same problem applies to notifying other microservices, updating search indexes, etc.


  How does it work? (The Solution: The Transactional Outbox Pattern)
  This is the clever part. We solve the problem by using the database's own reliability (transactions).
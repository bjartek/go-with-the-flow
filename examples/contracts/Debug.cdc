pub contract Debug {

   pub event Log(msg: String)

   pub fun log(_ msg: String) : String {
       emit Log(msg: msg)
       return msg
   }

}
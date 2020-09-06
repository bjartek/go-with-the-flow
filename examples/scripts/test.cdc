// test script to ensure code is running
pub fun main(account: Address): String {
    return getAccount(account).address.toString()
}
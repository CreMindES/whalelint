package whalelint

import com.google.gson.GsonBuilder
import com.google.gson.annotations.SerializedName
import com.google.gson.reflect.TypeToken

class ValidationResult {

    class Issue {
        // @JvmField
        @SerializedName("Rule")
        var rule: Rule? = null

        @SerializedName("IsViolated")
        var isViolated = false

        @JvmField
        @SerializedName("LocationRange")
        var locationRange: LocationRange? = null

        // @JvmField
        @SerializedName("Message")
        var message: String? = null

        override fun toString(): String {
            return "{" +
                        "Description=" + rule!!.definition + "," +
                        "IsViolated=" + isViolated + "," +
                        "Message=" + message + "," +
                        "LocationRange: {" +
                            "Start: {" +
                                "LineNumber:" + locationRange!!.start!!.lineNumber +
                                "CharNumber:" + locationRange!!.start!!.charNumber +
                        "}," +
                            "End: {" +
                                "LineNumber:" + locationRange!!.start!!.lineNumber +
                               "CharNumber:" + locationRange!!.start!!.charNumber +
                            "}" +
                        "}" +
                    "}"
        }
    }

    inner class Rule {
        @SerializedName("ID")
        val ruleId: String? = null

        @SerializedName("Definition")
        val definition: String? = null

        @SerializedName("Description")
        val description: String? = null

        @SerializedName("Severity")
        val severity: String? = null
    }

    inner class Location {
        @SerializedName("LineNumber")
        val lineNumber = 0

        @SerializedName("CharNumber")
        val charNumber = 0
    }

    inner class LocationRange {
        @SerializedName("Start")
        val start: Location? = null

        @SerializedName("End")
        val end: Location? = null
    }

    companion object {
        fun parse(json: String): List<Issue> {
            val g = GsonBuilder().create()
            val listType = object : TypeToken<List<Issue?>?>() {}.type

            return g.fromJson(json, listType)
        }
    }
}

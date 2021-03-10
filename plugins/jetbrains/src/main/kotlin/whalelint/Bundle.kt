package whalelint

import com.intellij.AbstractBundle
import org.jetbrains.annotations.PropertyKey
import java.util.*

object Bundle {
    private const val BUNDLE = "whalelint.Bundle"
    private val bundle = ResourceBundle.getBundle(BUNDLE)

    fun message(@PropertyKey(resourceBundle = BUNDLE) key: String, vararg params: Any): String {
        return AbstractBundle.message(bundle, key, *params)
    }
}

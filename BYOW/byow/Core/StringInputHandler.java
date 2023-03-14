package byow.Core;


public class StringInputHandler {
    private String Input;
    private int index;
    public StringInputHandler(String str) {
        Input = str;
        index = 0;
    }
    public char nextKeyTyped() {
        char returnChar = Input.charAt(index);
        index += 1;
        return returnChar;
    }

    public boolean hasNextKeyTyped() {
        return index < Input.length();
    }

}

package byow.Core;
import edu.princeton.cs.algs4.In;
import edu.princeton.cs.algs4.Out;

public class SaveHandler {
    public final static String filename = "Save.txt";

    public static void SaveProgress(String progress) {
        Out out = new Out(filename);
        out.println(progress);
    }

    public static String LoadProgress() {
        In in = new In(filename);
        String a = in.readAll();
        return a.substring(0, a.length() - 4);
    }

    public static void main(String[] args) {
        //System.out.println("Writing to " + filename);
        //Out out = new Out(filename);
        //out.println("This text will appear in the file!");
        //out.println("Your lucky number is " + System.currentTimeMillis() % 100);
        String a = LoadProgress();
        System.out.println(a);
    }

}

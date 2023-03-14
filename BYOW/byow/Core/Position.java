package byow.Core;

import org.junit.jupiter.api.ClassOrderer;
import java.util.Random;

public class Position {
    private int xPos;
    private int yPos;
    public Position(int x, int y) {
        xPos = x;
        yPos = y;
    }

    public int[] getCoor() {
        int[] arr = new int[2];
        arr[0] = xPos;
        arr[1] = yPos;
        return arr;
    }
    public int getxPos() {
        return xPos;
    }
    public int getyPos() {
        return yPos;
    }
    public Position changeX(int amount) {
        return new Position(this.getxPos()+amount, this.getyPos());
    }
    public Position changeY(int amount) {
        return new Position(this.getxPos(), this.getyPos()+amount);
    }
    public static Position RandomCoor(Random RANDOM, int xlower, int xupper, int ylower, int yupper) {
        return new Position(RANDOM.nextInt(xlower, xupper), RANDOM.nextInt(ylower, yupper));
    }
    public static Position RandomCoor(Random RANDOM, int xupper, int yupper) {
        return new Position(RANDOM.nextInt(xupper), RANDOM.nextInt(yupper));
    }

}

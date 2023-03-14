package byow.Core;

import byow.TileEngine.TETile;
import byow.TileEngine.Tileset;

import java.awt.*;

public class Player implements GameObject{
    private int health = 100;
    private Position coor;
    private TETile representation = Tileset.AVATAR;
    private Color color = Color.YELLOW;


    public Player(Position start) {
        coor = start;
    }

    public int GetHealth() {
        return health;
    }

    public Position GetPos() {
        return coor;
    }
    public void MoveUp() {
        coor = coor.changeY(1);
    }
    public void MoveDown() {
        coor = coor.changeY(-1);
    }
    public void MoveLeft() {
        coor = coor.changeX(-1);
    }
    public void MoveRight() {
        coor = coor.changeX(1);
    }
    public TETile GetRepresentation() {
        return representation;
    }
    public void ChangeRepresentation(TETile changed) {
        representation = changed;
    }
    public Color GetColor() {
        return color;
    }
    public void ChangeColor(Color changed) {
        color = changed;
    }
}

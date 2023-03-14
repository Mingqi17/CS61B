package byow.Core;

import byow.TileEngine.TERenderer;
import byow.TileEngine.TETile;
import byow.TileEngine.Tileset;
import edu.princeton.cs.algs4.StdDraw;

import java.sql.Array;
import java.util.ArrayList;
import java.util.Random;
import java.util.TooManyListenersException;

public class Hallway {

    private static final int WIDTH = 50;
    private static final int HEIGHT = 50;

    private static final long SEED = 11821;

    private static final Random RANDOM = new Random(SEED);

    /* Choose a place by random, given an input of a room or wall of the room */
    public static void chooseAPlace () {

        /* return the place choosen. Maybe a corr */
    }

    /* Input: a WALL, and convert it to FLOOR */
    public static void wallToFloor(TETile[][] tiles, Position pos) {
        /* What if a place is not a WALL? How to tell what the current place is? WALL, FLOOR, or other things? */
        tiles[pos.getxPos()][pos.getyPos()] = Tileset.FLOOR;
    }

    public static void fillNothing(TETile[][] tiles) {
        int height = tiles[0].length;
        int width = tiles.length;
        for (int x = 0; x < width; x += 1) {
            for (int y = 0; y < height; y += 1) {
                tiles[x][y] = Tileset.NOTHING;
            }
        }
    }

    /* connect: find the coor of intersection turns it into FLOOR and connect both way  */

    private static void connectTwoLower (TETile[][] tiles, Position pos1, Position pos2, TETile type) {
        int startX, endX;
        int startY, endY;
        int lineoflow;
        if (pos1.getxPos() < pos2.getxPos()) {
            startX = pos1.getxPos();
            endX = pos2.getxPos();
        } else {
            startX = pos2.getxPos();
            endX = pos1.getxPos();
        }
        if (pos1.getyPos() < pos2.getyPos()) {
            startY = pos1.getyPos();
            endY = pos2.getyPos();
            lineoflow = pos2.getxPos();
        } else {
            startY = pos2.getyPos();
            endY = pos1.getyPos();
            lineoflow = pos1.getxPos();
        }


        /* Fill in the FLOOR */
        for (int yPos = startY; yPos <= endY; yPos++) {
            //if (tiles[startX][yPos] == Tileset.FLOOR) {
            //endX = startX;
            //endY = yPos;
            //  break;
            // }
            if (type == Tileset.FLOOR || tiles[lineoflow][yPos] == Tileset.NOTHING) {
                tiles[lineoflow][yPos] = type;
            }

        }

        for (int xPos = startX; xPos <= endX; xPos++) {
            // if (tiles[xPos][endY] == Tileset.FLOOR) {
            //   endX = xPos;
            //   break;
            //   }

            if (type == Tileset.FLOOR || tiles[xPos][startY] == Tileset.NOTHING) {
                tiles[xPos][startY] = type;
            }
        }

    }

    private static void connectTwoHigher (TETile[][] tiles, Position pos1, Position pos2, TETile type) {
        int startX, endX;
        int startY, endY;
        int lineoflow;
        if (pos1.getxPos() < pos2.getxPos()) {
            startX = pos1.getxPos();
            endX = pos2.getxPos();
        } else {
            startX = pos2.getxPos();
            endX = pos1.getxPos();
        }
        if (pos1.getyPos() < pos2.getyPos()) {
            startY = pos1.getyPos();
            endY = pos2.getyPos();
            lineoflow = pos1.getxPos();
        } else {
            startY = pos2.getyPos();
            endY = pos1.getyPos();
            lineoflow = pos2.getxPos();
        }

        /* Fill in the FLOOR */
        for (int yPos = startY; yPos <= endY; yPos++) {
            //if (tiles[startX][yPos] == Tileset.FLOOR) {
                //endX = startX;
                //endY = yPos;
              //  break;
           // }
            if (type == Tileset.FLOOR || tiles[lineoflow][yPos] == Tileset.NOTHING) {
                tiles[lineoflow][yPos] = type;
            }

        }
        for (int xPos = startX; xPos <= endX; xPos++) {
           // if (tiles[xPos][endY] == Tileset.FLOOR) {
             //   endX = xPos;
             //   break;
         //   }

            if (type == Tileset.FLOOR || tiles[xPos][endY] == Tileset.NOTHING) {
                tiles[xPos][endY] = type;
            }
        }

    }
    public static void drawhallwayHigher(TETile[][] tiles, Position pos1, Position pos2) {
        if (pos1.getxPos() == pos2.getxPos()) {
            drawVerStraight(tiles, pos1, pos2);
            return;
        } else if (pos1.getyPos() == pos2.getyPos()) {
            drawHorStraight(tiles, pos1, pos2);
            return;
        }
        connectTwoHigher(tiles, pos1, pos2, Tileset.FLOOR);
        int xMax= Integer.max(pos1.getxPos(), pos2.getxPos());
        int yMax= Integer.max(pos1.getyPos(), pos2.getyPos());
        int xMin= Integer.min(pos1.getxPos(), pos2.getxPos());
        int yMin= Integer.min(pos1.getyPos(), pos2.getyPos());
        if ((pos1.getxPos() > pos2.getxPos() && pos1.getyPos() > pos2.getyPos()) ||
                (pos1.getxPos() < pos2.getxPos() && pos1.getyPos() < pos2.getyPos())) {
            connectTwoHigher(tiles, new Position(xMax, yMax + 1), new Position(xMin - 1, yMin), Tileset.WALL);
            connectTwoHigher(tiles, new Position(xMax, yMax - 1), new Position(xMin + 1, yMin), Tileset.WALL);
        } else {
            connectTwoHigher(tiles, new Position(xMax + 1, yMin), new Position(xMin, yMax + 1), Tileset.WALL);
            connectTwoHigher(tiles, new Position(xMax - 1, yMin), new Position(xMin, yMax - 1), Tileset.WALL);
        }

    }
    public static void drawhallwayLower(TETile[][] tiles, Position pos1, Position pos2) {
        if (pos1.getxPos() == pos2.getxPos()) {
            drawVerStraight(tiles, pos1, pos2);
            return;
        } else if (pos1.getyPos() == pos2.getyPos()) {
            drawHorStraight(tiles, pos1, pos2);
            return;
        }

        connectTwoLower(tiles, pos1, pos2, Tileset.FLOOR);
        int xMax= Integer.max(pos1.getxPos(), pos2.getxPos());
        int yMax= Integer.max(pos1.getyPos(), pos2.getyPos());
        int xMin= Integer.min(pos1.getxPos(), pos2.getxPos());
        int yMin= Integer.min(pos1.getyPos(), pos2.getyPos());
        if ((pos1.getxPos() > pos2.getxPos() && pos1.getyPos() > pos2.getyPos()) ||
                (pos1.getxPos() < pos2.getxPos() && pos1.getyPos() < pos2.getyPos())) {
            connectTwoLower(tiles, new Position(xMin, yMin - 1), new Position(xMax + 1, yMax), Tileset.WALL);
            connectTwoLower(tiles, new Position(xMin, yMin + 1), new Position(xMax - 1, yMax), Tileset.WALL);
        } else {
            connectTwoLower(tiles, new Position(xMax, yMin + 1), new Position(xMin + 1, yMax), Tileset.WALL);
            connectTwoLower(tiles, new Position(xMax, yMin - 1), new Position(xMin - 1, yMax), Tileset.WALL);
        }

    }
    private static void drawVerStraight(TETile[][] tiles, Position pos1, Position pos2) {
        int X = pos1.getxPos();
        int firstY = pos1.getyPos();
        int secondY = pos2.getyPos();
        int startY, endY;
        if (firstY > secondY) {
            startY = secondY;
            endY = firstY;
        } else {
            startY = firstY;
            endY = secondY;
        }
        connectVerStraight(tiles, X, startY, endY, Tileset.FLOOR);
        connectVerStraight(tiles, X + 1, startY, endY, Tileset.WALL);
        connectVerStraight(tiles, X - 1, startY, endY, Tileset.WALL);
    }
    private static void drawHorStraight(TETile[][] tiles, Position pos1, Position pos2) {
        int Y = pos1.getyPos();
        int firstX = pos1.getxPos();
        int secondX = pos2.getxPos();
        int startX, endX;
        if (firstX > secondX) {
            startX = secondX;
            endX = firstX;
        } else {
            startX= firstX;
            endX = secondX;
        }
        connectHorStraight(tiles, Y, startX, endX, Tileset.FLOOR);
        connectHorStraight(tiles, Y + 1, startX, endX, Tileset.WALL);
        connectHorStraight(tiles, Y - 1, startX, endX, Tileset.WALL);
    }
    private static void connectVerStraight(TETile[][] tiles, int X, int startY, int endY, TETile type) {
        for (int i = startY; i <= endY; i ++) {
            if (type == Tileset.FLOOR || tiles[X][i] == Tileset.NOTHING) {
                tiles[X][i] = type;
            }
        }
    }
    private static void connectHorStraight(TETile[][] tiles, int Y, int startX, int endX, TETile type) {
        for (int i = startX; i <= endX; i ++) {
            if (type == Tileset.FLOOR || tiles[i][Y] == Tileset.NOTHING) {
                tiles[i][Y] = type;
            }
        }
    }

    public static void main(String[] args) {
        TERenderer ter = new TERenderer();
        ter.initialize(WIDTH, HEIGHT);

        /* See which way to fill in is better */

        TETile[][] randomTiles = new TETile[WIDTH][HEIGHT];
        fillNothing(randomTiles);
        Position pos1 = Position.RandomCoor(RANDOM, 11, 20, 10, 20);
        Position pos2 = Position.RandomCoor(RANDOM, 1, 10, 1, 10);
        Position pos3 = new Position(28,20);
        Position pos4 = new Position(20,20);
        connectTwoLower(randomTiles, pos1, pos2, Tileset.FLOOR);
        drawhallwayLower(randomTiles, pos3, pos4);
        //drawhallway(randomTiles, pos3, pos4);
        //onnectTwoLower(randomTiles, pos1, pos2);

        ter.renderFrame(randomTiles);
        StdDraw.show();
    }
}

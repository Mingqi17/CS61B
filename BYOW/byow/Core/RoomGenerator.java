package byow.Core;
import byow.TileEngine.TERenderer;
import byow.TileEngine.TETile;
import byow.TileEngine.Tileset;
import edu.princeton.cs.algs4.StdDraw;

import java.awt.*;
import java.util.*;
import java.util.List;

public class RoomGenerator {
    private static final int WIDTH = 50;
    private static final int HEIGHT = 50;

    private final long SEED;
    private final Random RANDOM;
    private static final int HORDISTANCE = WIDTH/5;
    private static final int VERDISTANCE = HEIGHT/5;
    private static final int ROOMSIZE = 8;


    public RoomGenerator(int seed) {
        SEED = seed;
        RANDOM = new Random(SEED);
    }

    /** Check whether the room can be generated at this place
     * @param tiles the current world
     * @param coor the coordinate to generate the room
     * @param width the width of the prospective room
     * @param length the length of the prospective room
     * @return boolean whether it is possible
     */

    private static boolean CheckAval(TETile[][] tiles, Position coor, int length, int width) {
        int x = coor.getxPos();
        int y = coor.getyPos();
        if (x + width >= WIDTH || y + length >= HEIGHT || y <= 0 || x <= 0) {
            return false;
        }
        for (int i = x - 1; i <= x + width; i += 1) {
            for (int j = y - 1; j <= y + length; j += 1) {
                if (tiles[i][j] != Tileset.NOTHING) {
                    return false;
                }
            }
        }
        return true;
    }

    /** Build a Room
     * @param tiles the current world
     * @param coor the coordinate to generate the room
     * @param width the width of the room
     * @param length the length of the room
     * @return boolean whether the room is successfully generated
     */
    private static boolean BuildRoom(TETile[][] tiles, Position coor, int width, int length) {
        int x = coor.getxPos();
        int y = coor.getyPos();
        if (CheckAval(tiles, coor, length, width)) {
            for (int i = x; i < x + width; i += 1) {
                for (int j = y; j < y + length; j += 1) {
                    tiles[i][j] = Tileset.FLOOR;
                }
            }
            for (int i = x - 1; i <= x + width; i += 1) {
                tiles[i][y - 1] = Tileset.WALL;
                tiles[i][y + length] = Tileset.WALL;
            }
            for (int j = y - 1; j < y + length; j += 1) {
                tiles[x - 1][j] = Tileset.WALL;
                tiles[x + width][j] = Tileset.WALL;
            }
            return true;
        }
        return false;
    }

    /** Give another 5 possible positions for new rooms from the current coordinate
     * @param coor current coordinate
     * @return an ArrayList of new positions
     */
    private static ArrayList<List<Position>> NearestBuild(Position coor, Random RANDOM) {
        ArrayList<List<Position>> toBuild = new ArrayList<>();
        int x = coor.getxPos();
        int y = coor.getyPos();
        for (int i = 0; i < 5; i ++) {
            List<Position> curr = new ArrayList<>();
            curr.add(coor);
            curr.add(Position.RandomCoor(RANDOM, (-1 * HORDISTANCE) + x, HORDISTANCE + x,
                    (-1 * VERDISTANCE) + y, VERDISTANCE + y));
            toBuild.add(curr);
        }
        return toBuild;
    }

    /** Initiate random size for a room
     * @return a list of two elements: width and length
     */
    private static int[] RandomSize(Random RANDOM) {
        return new int[]{RANDOM.nextInt(1, ROOMSIZE), RANDOM.nextInt(1, ROOMSIZE)};
    }

    /** Generate the Room part of this world
     * @param tiles the current world
     */
    public void RoomGenerate(TETile[][] tiles) {
        int number = RANDOM.nextInt(40, 60);
        Position center = new Position(WIDTH/2, HEIGHT/2);
        Queue<List<Position>> toDoRoom = new LinkedList<>();
        ArrayList<Position[]> HallwaystoDo = new ArrayList<>();
        ArrayList<Position> start = new ArrayList<>();
        start.add(center);
        start.add(center);
        toDoRoom.add(start);
        while (number > 0 && toDoRoom.size() > 0) {
            List<Position> current = toDoRoom.poll();
            int[] size = RandomSize(RANDOM);
            if (BuildRoom(tiles, current.get(1), size[0], size[1])) {
                ArrayList<List<Position>> nextToDo = NearestBuild(current.get(1), RANDOM);
                Position currRan = Position.RandomCoor(RANDOM,current.get(1).getxPos(), current.get(1).getxPos() + size[0],
                        current.get(1).getyPos(), current.get(1).getyPos() + size[1]);
                Position[] connection = {currRan, current.get(0)};
                HallwaystoDo.add(connection);
                toDoRoom.addAll(nextToDo);
                number -= 1;
            }
        }
        HallwaysHandle(tiles, HallwaystoDo, RANDOM);

    }

    private static void HallwaysHandle(TETile[][] tiles, ArrayList<Position[]> toDo, Random RANDOM) {
        for (Position[] i : toDo) {
            int decision = RANDOM.nextInt(2);
            if (decision == 0) {
                Hallway.drawhallwayHigher(tiles, i[1], i[0]);
            } else {
                Hallway.drawhallwayLower(tiles, i[1], i[0]);
            }
        }
    }






}

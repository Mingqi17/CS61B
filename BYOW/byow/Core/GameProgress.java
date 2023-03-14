package byow.Core;
import byow.TileEngine.TETile;
import byow.TileEngine.Tileset;
import edu.princeton.cs.algs4.StdDraw;

import java.awt.*;
import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;
import java.util.Queue;

public class GameProgress {
    private static final int WIDTH = 50;
    private static final int HEIGHT = 50;
    private static final int SIGHT = 4;
    public boolean BLOCKED = false;
    public boolean ATTEMPT = false;
    public boolean Quit = false;

    private static Position GetArrowPosition() {
        int X = (int) StdDraw.mouseX();
        int Y = (int) StdDraw.mouseY();
        return new Position(X, Y);
    }

    private static String GetPositionType(TETile[][] world, Position coor, Player player) {
        int x = coor.getxPos();
        int y = coor.getyPos();
        if (x < 0 || x >= WIDTH || y < 0 || y >= HEIGHT) {
            return "nothing";
        }
        int px = player.GetPos().getxPos();
        int py = player.GetPos().getyPos();
        if (px == x & py == y) {
            return "you";
        }
        return world[x][y].description();

    }
    public static String GetArrowType(TETile[][] world, Player player) {
        return GetPositionType(world, GetArrowPosition(), player);
    }

    public TETile[][] AvartarInput(TETile[][] world, Player player, char Input) {
        Position playerpos = player.GetPos();
        switch (Input) {
            case 'w':
                if (CheckPossibleMove(world, playerpos.changeY(1))) {
                    world[playerpos.getxPos()][playerpos.getyPos()] = Tileset.FLOOR;
                    player.MoveUp();
                    playerpos = player.GetPos();
                    world[playerpos.getxPos()][playerpos.getyPos()] = player.GetRepresentation();
                }
                ATTEMPT = false;
                return world;
            case 's':
                if (CheckPossibleMove(world, playerpos.changeY(-1))) {
                    world[playerpos.getxPos()][playerpos.getyPos()] = Tileset.FLOOR;
                    player.MoveDown();
                    playerpos = player.GetPos();
                    world[playerpos.getxPos()][playerpos.getyPos()] = player.GetRepresentation();
                }
                ATTEMPT = false;
                return world;
            case 'a':
                if (CheckPossibleMove(world, playerpos.changeX(-1))) {
                    world[playerpos.getxPos()][playerpos.getyPos()] = Tileset.FLOOR;
                    player.MoveLeft();
                    playerpos = player.GetPos();
                    world[playerpos.getxPos()][playerpos.getyPos()] = player.GetRepresentation();
                }
                ATTEMPT = false;
                return world;
            case 'd':
                if (CheckPossibleMove(world, playerpos.changeX(1))) {
                    world[playerpos.getxPos()][playerpos.getyPos()] = Tileset.FLOOR;
                    player.MoveRight();
                    playerpos = player.GetPos();
                    world[playerpos.getxPos()][playerpos.getyPos()] = player.GetRepresentation();
                }
                ATTEMPT = false;
                return world;
            case 'f':
                BLOCKED = !BLOCKED;
                ATTEMPT = false;
                return world;
            case ':':
                ATTEMPT = true;
                return world;
            case 'q':
                if (ATTEMPT) {
                    Quit = true;
                }
            case 'Q':
                if (ATTEMPT) {
                    Quit = true;
                }
            default:
                ATTEMPT = false;
                return world;
        }
    }

    private static boolean CheckPossibleMove(TETile[][] world, Position newcoor) {
        int x = newcoor.getxPos();
        int y = newcoor.getyPos();
        return world[x][y] == Tileset.FLOOR;
    }

    public static String CustomPlayer(Player player) {
        String result = "";
        GameMenu.DecisionMenu();
        Color[] colorsets = {Color.blue, Color.CYAN, Color.GREEN, Color.MAGENTA, Color.ORANGE, Color.red, Color.YELLOW};
        TETile[] charssets = {Tileset.MOUNTAIN, Tileset.WATER, Tileset.AVATAR};
        boolean finished = false;
        while (!finished) {
            if (StdDraw.hasNextKeyTyped()) {
                char next = StdDraw.nextKeyTyped();
                result += next;
                if (next == 'N' | next == 'n') {
                    finished = true;
                } else if (next == 'Y' | next == 'y') {
                    int i = 0;
                    int n = 0;
                    GameMenu.CustomMenu(player);
                    while (!finished) {
                        if (StdDraw.hasNextKeyTyped()) {
                            char curr = StdDraw.nextKeyTyped();
                            result += curr;
                            switch (curr) {
                                case 'c':
                                    player.ChangeColor(colorsets[i % colorsets.length]);
                                    i += 1;
                                    break;
                                case 'x':
                                    player.ChangeRepresentation(charssets[n % charssets.length]);
                                    n += 1;
                                    break;
                                case 'z':
                                    finished = true;
                                    break;
                                default:
                                    finished = false;
                                    break;
                            }
                            GameMenu.CustomMenu(player);

                        }
                    }

                }

            }
        }
        return result;
    }







    public static void LimitSight(TETile[][] world, Player player) {
        StdDraw.clear(new Color(0, 0, 0));
        int px = player.GetPos().getxPos();
        int py = player.GetPos().getyPos();
        int numXTiles = world.length;
        int numYTiles = world[0].length;
        boolean[][] checked = new boolean[numXTiles][numYTiles];
        //StdDraw.clear(new Color(0, 0, 0));
        for (int x = 0; x < numXTiles; x += 1) {
            for (int y = 0; y < numYTiles; y += 1) {
                checked[x][y] = false;
            }
        }
        checked[px][py] = true;




        Queue<DisPos> toDoPos = new LinkedList<>();
        toDoPos.add(new DisPos(player.GetPos(), 0));
        while (toDoPos.size() > 0) {
            DisPos curr = toDoPos.poll();
            ArrayList<DisPos> next = CheckNeighbor(checked, curr, world, player);
            if (next.size() > 0) {
                toDoPos.addAll(next);
            }
        }

        for (int x = 0; x < numXTiles; x += 1) {
            for (int y = 0; y < numYTiles; y += 1) {
                if (checked[x][y]) {
                    world[x][y].draw(x, y);
                }
            }
        }

    }

    private static ArrayList<DisPos> CheckNeighbor(boolean[][] checked, DisPos curr, TETile[][] world, Player player) {
        ArrayList<DisPos> nextToDo = new ArrayList<>();
        int x = curr.coor.getxPos();
        int y = curr.coor.getyPos();
        int px = player.GetPos().getxPos();
        int py = player.GetPos().getyPos();
        int Dis = curr.Dis;
        if (!checked[x + 1][y] & world[x + 1][y] != Tileset.WALL & Dis < SIGHT) {
            nextToDo.add(new DisPos(new Position(x + 1, y), Dis + 1));
        }
        if (!checked[x - 1][y] & world[x - 1][y] != Tileset.WALL & Dis < SIGHT) {
            nextToDo.add(new DisPos(new Position(x - 1, y), Dis + 1));
        }
        if (!checked[x][y + 1] & world[x][y + 1] != Tileset.WALL & Dis < SIGHT) {
            nextToDo.add(new DisPos(new Position(x, y + 1), Dis + 1));
        }
        if (!checked[x][y - 1] & world[x][y - 1] != Tileset.WALL & Dis < SIGHT) {
            nextToDo.add(new DisPos(new Position(x, y - 1), Dis + 1));
        }
        if (!checked[x + 1][y + 1] & world[x + 1][y + 1] != Tileset.WALL & Dis + 1 < SIGHT) {
            nextToDo.add(new DisPos(new Position(x + 1, y + 1), Dis + 2));
        }
        if (!checked[x - 1][y + 1] & world[x - 1][y + 1] != Tileset.WALL & Dis + 1 < SIGHT) {
            nextToDo.add(new DisPos(new Position(x - 1, y + 1), Dis + 2));
        }
        if (!checked[x - 1][y - 1] & world[x - 1][y - 1] != Tileset.WALL & Dis + 1 < SIGHT) {
            nextToDo.add(new DisPos(new Position(x - 1, y - 1), Dis + 2));
        }
        if (!checked[x + 1][y - 1] & world[x + 1][y - 1] != Tileset.WALL & Dis + 1 < SIGHT) {
            nextToDo.add(new DisPos(new Position(x + 1, y - 1), Dis + 2));
        }
        if (checked[x][y]) {
            checked[x + 1][y] = true;
            checked[x - 1][y] = true;
            checked[x][y + 1] = true;
            checked[x][y - 1] = true;
        }
        if (checked[x][y] & Dis + 2 <= SIGHT) {
            checked[x + 1][y + 1] = true;
            checked[x - 1][y - 1] = true;
            checked[x - 1][y + 1] = true;
            checked[x + 1][y - 1] = true;
        }
        return nextToDo;
    }


}

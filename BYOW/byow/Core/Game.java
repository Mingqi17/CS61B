package byow.Core;

import byow.TileEngine.TERenderer;
import byow.TileEngine.TETile;
import byow.TileEngine.Tileset;
import edu.princeton.cs.algs4.In;
import edu.princeton.cs.algs4.ST;
import edu.princeton.cs.algs4.StdDraw;

import java.awt.*;
import java.io.FileWriter;
import java.net.Inet4Address;
import java.util.Map;
import java.util.logging.Handler;

public class Game {
    private int SEED;
    private boolean GameOver = false;
    private boolean GameStart = false;
    private static final int WIDTH = 50;
    private static final int HEIGHT = 50;
    private String SavedGame = SaveHandler.LoadProgress();
    public String CurrentProgress = "";
    private final static String filename = "Save.txt";
    private Player avatar = new Player(new Position(WIDTH / 2, HEIGHT / 2));

    public Game() {
    }
    public void EndGame() {
        GameOver = true;
    }
    public void startGame() {
        GameProgress trivial = new GameProgress();
        TERenderer ter = new TERenderer();
        ter.initialize(WIDTH, HEIGHT);
        String seed = InputSeed();
        TETile[][] world;
        if (seed != null) {
            SEED = Integer.parseInt(seed);
            world = MapGeneration.GenerateWorld(SEED, ter);
            world[WIDTH / 2][HEIGHT / 2] = avatar.GetRepresentation();
        } else {
            //System.out.println(SavedGame);
            CurrentProgress += SavedGame;
            world = InteractwithString(SavedGame);
        }


        while (!GameOver) {
            Font font = new Font("Monaco", Font.BOLD, 14);
            StdDraw.setFont(font);
            if (StdDraw.hasNextKeyTyped()) {
                char next = StdDraw.nextKeyTyped();
                CurrentProgress += next;
                world = trivial.AvartarInput(world, avatar, next);
            }

            if (trivial.Quit) {
                //Save Game
                SaveHandler.SaveProgress(CurrentProgress);
                EndGame();
                System.exit(0);
            }
            if (!trivial.BLOCKED) {
                ter.renderFrame(world);
            } else {
                GameProgress.LimitSight(world, avatar);
            }
            world[avatar.GetPos().getxPos()][avatar.GetPos().getyPos()].draw(avatar.GetPos().getxPos(),
                    avatar.GetPos().getyPos(), avatar.GetColor());
            HUD.DisplayHUD(world, avatar);
            StdDraw.show();
        }

    }

    public String InputSeed() {
        while (!GameStart) {
            GameMenu.MainMenu();
            if (StdDraw.hasNextKeyTyped()) {
                char next = StdDraw.nextKeyTyped();
                CurrentProgress += next;
                if (next == 'P' | next == 'p') {
                    CurrentProgress += GameProgress.CustomPlayer(avatar);
                }
                if (next == 'Q' | next == 'q') {
                    EndGame();
                    System.exit(0);
                }
                if (next == 'L' | next == 'l') {
                    break;
                }
                if (next == 'N' | next == 'n') {
                    GameMenu.SeedMenu("");
                    String result = "";
                    while (!GameStart) {
                        if (StdDraw.hasNextKeyTyped()) {
                            char curr = StdDraw.nextKeyTyped();
                            CurrentProgress += curr;
                            if (curr == 'S' | curr == 's') {
                                GameStart = true;
                                return result;
                            }
                            result += curr;
                            GameMenu.SeedMenu(result);
                        }
                    }
                }
            }
        }
        return null;
    }

    public TETile[][] InteractwithString(String Input) {
        StringInputHandler handler = new StringInputHandler(Input);
        GameProgress trivial = new GameProgress();
        String result = "";
        while (!GameStart) {
            char curr = handler.nextKeyTyped();
            if (curr == 'P' | curr == 'p') {
                Color[] colorsets = {Color.blue, Color.CYAN, Color.GREEN, Color.MAGENTA, Color.ORANGE, Color.red, Color.YELLOW};
                TETile[] charssets = {Tileset.MOUNTAIN, Tileset.WATER, Tileset.AVATAR};
                boolean finished = false;
                while (!finished) {
                    if (handler.hasNextKeyTyped()) {
                        char next = handler.nextKeyTyped();
                        if (next == 'N' | next == 'n') {
                            finished = true;
                            break;
                        } else if (next == 'Y' | next == 'y') {
                            int i = 0;
                            int n = 0;
                            while (!finished) {
                                if (handler.hasNextKeyTyped()) {
                                    char first = handler.nextKeyTyped();
                                    switch (first) {
                                        case 'c':
                                            avatar.ChangeColor(colorsets[i % colorsets.length]);
                                            i += 1;
                                            break;
                                        case 'x':
                                            avatar.ChangeRepresentation(charssets[n % charssets.length]);
                                            n += 1;
                                            break;
                                        case 'z':
                                            finished = true;
                                            break;
                                        default:
                                            finished = false;
                                            break;
                                    }
                                }
                            }

                        }

                    }
                }

            }

            if (curr == 'N' | curr == 'n') {
                while (!GameStart) {
                    if (handler.hasNextKeyTyped()) {
                        char next = handler.nextKeyTyped();
                        if (next == 'S' | next == 's') {
                            GameStart = true;
                            break;
                        }
                        result += next;
                    }
                }
            }
        }

        SEED = Integer.parseInt(result);
        TETile[][] world = MapGeneration.GenerateWorld(SEED);
        if (!handler.hasNextKeyTyped()) {
            world[WIDTH / 2][HEIGHT / 2] = avatar.GetRepresentation();
        }
        while (handler.hasNextKeyTyped()) {
            world = trivial.AvartarInput(world, avatar, handler.nextKeyTyped());
        }

        //TERenderer ter = new TERenderer();
        // ter.initialize(WIDTH, HEIGHT);
        //ter.renderFrame(world);
        //StdDraw.show();



        return world;
    }




    public static void main(String[] args) {
        Game newGame = new Game();
        newGame.startGame();
    }

}

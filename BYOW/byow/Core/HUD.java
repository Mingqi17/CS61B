package byow.Core;

import byow.TileEngine.TETile;
import edu.princeton.cs.algs4.StdDraw;

import java.awt.*;

public class HUD {
    private static final int WIDTH = 50;
    private static final int HEIGHT = 50;

    private static void DisplayString(String hud) {
        StdDraw.setPenColor(Color.WHITE);
        Font fontBig = new Font("Monaco", Font.BOLD, 30);
        StdDraw.setFont(fontBig);
        StdDraw.textLeft( WIDTH / 20,HEIGHT * 9 / 10, hud);
    }

    public static void DisplayHUD(TETile[][] world, Player player) {
        DisplayString(GameProgress.GetArrowType(world, player));
    }
}

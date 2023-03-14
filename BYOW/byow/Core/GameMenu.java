package byow.Core;

import edu.princeton.cs.algs4.StdDraw;

import java.awt.*;
import java.util.Random;

import static byow.Core.SaveHandler.LoadProgress;

public class GameMenu {
    private static final int WIDTH = 50;
    private static final int HEIGHT = 50;

    public static void MainMenu() {
        StdDraw.clear(new Color(0, 0, 0));
        StdDraw.setPenColor(Color.WHITE);
        Font fontBig = new Font("Monaco", Font.BOLD, 50);
        StdDraw.setFont(fontBig);
        StdDraw.text( WIDTH / 2,HEIGHT * 2 / 3, "Game");
        Font fontSmall = new Font("Monaco", Font.BOLD, 30);
        StdDraw.setFont(fontSmall);
        StdDraw.text( WIDTH / 2,HEIGHT * 6 / 12, "Start New Game(N)");
        StdDraw.text( WIDTH / 2,HEIGHT * 5 / 12, "Load Recent Game(L)");
        StdDraw.text( WIDTH / 2,HEIGHT * 4 / 12, "Change Appearance(P)");
        StdDraw.text( WIDTH / 2,HEIGHT * 3 / 12, "Quit Game(Q)");
        StdDraw.show();
    }

    public static void SeedMenu(String curr) {
        StdDraw.clear(new Color(0, 0, 0));
        StdDraw.setPenColor(Color.WHITE);
        Font fontBig = new Font("Monaco", Font.BOLD, 30);
        StdDraw.setFont(fontBig);
        StdDraw.text( WIDTH / 2,HEIGHT * 2 / 3, "Please Enter Your Seed, Press S to proceed");
        StdDraw.text( WIDTH / 2,HEIGHT / 2, curr);
        StdDraw.show();
    }
    public static void CustomMenu(Player player) {
        Color color = player.GetColor();
        char rep = player.GetRepresentation().character();
        StdDraw.clear(new Color(0, 0, 0));
        StdDraw.setPenColor(Color.WHITE);
        Font fontBig = new Font("Monaco", Font.BOLD, 30);
        StdDraw.setFont(fontBig);
        StdDraw.text( WIDTH / 2,HEIGHT * 2 / 3, "Press C to get change Color");
        StdDraw.text( WIDTH / 2,HEIGHT / 2, "Press X to change character");
        StdDraw.text( WIDTH / 2,HEIGHT / 3, "Press Z to confirm");
        StdDraw.setPenColor(color);
        Font character = new Font("Monaco", Font.BOLD, 60);
        StdDraw.setFont(character);
        StdDraw.text( WIDTH / 2,HEIGHT / 6, Character.toString(rep));
        StdDraw.show();
    }
    public static void DecisionMenu() {
        StdDraw.clear(new Color(0, 0, 0));
        StdDraw.setPenColor(Color.WHITE);
        Font fontBig = new Font("Monaco", Font.BOLD, 30);
        StdDraw.setFont(fontBig);
        StdDraw.text( WIDTH / 2,HEIGHT * 2 / 3, "Do you wish to change the appearance?");
        StdDraw.text( WIDTH / 2,HEIGHT / 2, "Yes(Y) and No(N)");
        StdDraw.show();
    }



}

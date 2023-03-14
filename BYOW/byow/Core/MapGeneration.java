package byow.Core;

import byow.TileEngine.TERenderer;
import byow.TileEngine.TETile;
import byow.TileEngine.Tileset;

import java.util.Random;

public class MapGeneration {
    private static final int WIDTH = 50;
    private static final int HEIGHT = 50;

    private static final long SEED = 2873123;
    private static final Random RANDOM = new Random(SEED);

    /** Initiate the world with all blocks to be nothing
     * @param tiles the current world
     */
    public static void fillNothing(TETile[][] tiles) {
        int height = tiles[0].length;
        int width = tiles.length;
        for (int x = 0; x < width; x += 1) {
            for (int y = 0; y < height; y += 1) {
                tiles[x][y] = Tileset.NOTHING;
            }
        }
    }

    public static TETile[][] GenerateWorld(int seed, TERenderer ter) {
        TETile[][] randomTiles = new TETile[WIDTH][HEIGHT];
        fillNothing(randomTiles);
        RoomGenerator newworld = new RoomGenerator(seed);
        newworld.RoomGenerate(randomTiles);
        ter.renderFrame(randomTiles);
        return randomTiles;
    }
    public static TETile[][] GenerateWorld(int seed) {
        TETile[][] randomTiles = new TETile[WIDTH][HEIGHT];
        fillNothing(randomTiles);
        RoomGenerator newworld = new RoomGenerator(seed);
        newworld.RoomGenerate(randomTiles);
        return randomTiles;
    }

}

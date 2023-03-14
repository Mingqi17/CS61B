package byow.Core;

public interface GameObject {
    public Position GetPos();
    public int GetHealth();
    public void MoveUp();
    public void MoveDown();
    public void MoveLeft();
    public void MoveRight();
}

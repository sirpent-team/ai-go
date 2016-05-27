#define GRID_SIZE 300

struct snake {
  int snake_id;
  struct snake_segment * head;
  /* also put some network stuff here */
};

struct snake_segment {
  struct snake * snake_id;
  struct snake_segment * forward; // nullptr if head
  struct snake_segment * rearward; // nullptr if tip of tail
  int x_coord;
  int y_coord;
};

static struct snake_segment * grid0 [GRID_SIZE][GRID_SIZE];
static struct snake_segment * grid1 [GRID_SIZE][GRID_SIZE];

typedef enum {
  NORTH,
  NORTHEAST,
  SOUTHEAST,
  SOUTH,
  SOUTHWEST,
  NORTHWEST
} direction;

direction ask_direction(struct snake snake) {
  /* TODO actually ask */
  return SOUTH;
}

void update_world() {
  int x; int y;
  for (x = 0; x < GRID_SIZE; x++) {
    for (y = 0; y < GRID_SIZE; y++) {
      if (grid0[x][y]) {
	if (grid0[x][y]->forward) {
	  // we are in the tail of a snake
	  struct snake_segment forward = *grid0[x][y]->forward;
	  grid1[forward.x_coord][forward.y_coord] = grid0[x][y];
	}
	else {
	  // we are in the head of a snake
	  direction which_way = ask_direction(*grid0[x][y]->snake_id);
	  int new_x = x, new_y = y;
	  switch (which_way) {
	  case NORTH:
	    new_y -= 1;
	    break;
	  case NORTHEAST:
	    new_x += 1;
	    new_y -= 1;
	    break;
	  case SOUTHEAST:
	    new_x += 1;
	    break;
	  case SOUTH:
	    new_y += 1;
	    break;
	  case SOUTHWEST:
	    new_x -= 1;
	    new_y += 1;
	    break;
	  case NORTHWEST:
	    new_x -= 1;
	    break;
	  }
	  /* CHECK FOR COLLISIONS */
	  /* not actually sure the best way to do this??? */
	  /* make a second grid just for heads and merge later? */
	  grid1[new_x][new_y] = grid0[x][y];
	}
      }
    }
  }
  for (x = 0; x < GRID_SIZE; x++) {
    for (y = 0; y < GRID_SIZE; y++) {
      grid0[x][y] = grid1[x][y];
      if (grid1[x][y]) {
	grid1[x][y]->x_coord = x;
	grid1[x][y]->y_coord = y;
      }
      grid1[x][y] = 0;
    }
  }
  // TODO: copy grid1 into grid2
}

int main() {
  return 0;
}

import curses
import random

import numpy as np


class Game:
    symbols = ["A", "B", "C", "D", "E", "F"]
    empty_symbol = " "

    def __init__(self, stdscr, colors):
        self.stdscr = stdscr
        self.colors = colors
        self.grid = np.fromfunction(np.vectorize(lambda i, j: Game.get_random_symbol()), [10] * 2)
        # self.grid = [[self.get_random_symbol() for _ in range(Game.grid_size)] for _ in range(Game.grid_size)]

    @staticmethod
    def get_random_symbol():
        index = random.randrange(len(Game.symbols))
        return Game.symbols[index]

    def play(self):
        self.refresh_grid()

        while True:
            selection_valid = False
            while not selection_valid:
                selected_point1, selected_point2 = self.select_points_to_swap()
                prev_grid = self.grid.copy()
                self.swap(selected_point1, selected_point2)

                selection_valid = self.find_match() is not None
                if selection_valid:
                    self.draw(f"Swapped {self.grid[selected_point1[0], selected_point1[1]]} {selected_point1} "
                              f"and {self.grid[selected_point2[0], selected_point2[1]]} {selected_point2}\n"
                              f"Press any key to continue...", None)
                    self.stdscr.getch()
                else:
                    self.grid = prev_grid
                    self.draw(hint="Swap would not result in a match; please try again\nPress any key to continue...")
                    self.stdscr.getch()

            self.refresh_grid()

    def select_points_to_swap(self):
        def generate_hint(selected_position, selected_position_index):
            return (
                f"Select two symbols to swap...\n"
                f"Selecting symbol #{selected_position_index + 1}\n"
                f"Press arrow keys (← ↑ → ↓) to move selection; press enter to continue\n"
                f"{self.grid[selected_position[0], selected_position[1]]} {selected_position}"
            )

        def get_available_neighbor(selected_position1):
            if selected_position1[0] == 0:
                if selected_position1[1] == 0:
                    return selected_position1 + np.array([0, 1])
                else:
                    return selected_position1 + np.array([-1, 0])
            return selected_position1 + np.array([-1, 0])

        # Initialise selected position to center
        selected_position1 = (np.array(self.grid.shape) - 1) // 2
        self.draw(generate_hint(selected_position1, 0), [selected_position1])
        selecting = True
        while selecting:
            key = self.stdscr.getch()

            # If key is an arrow key then move the selected position
            # If key is enter then set selected position
            match key:
                case curses.KEY_UP:
                    selected_position1[0] -= 1
                case curses.KEY_DOWN:
                    selected_position1[0] += 1
                case curses.KEY_LEFT:
                    selected_position1[1] -= 1
                case curses.KEY_RIGHT:
                    selected_position1[1] += 1
                case curses.KEY_ENTER | 10 | 13:
                    selecting = False

            # Wrap selected position
            selected_position1 = selected_position1 % self.grid.shape

            self.draw(generate_hint(selected_position1, 0), [selected_position1])

        selected_position2 = get_available_neighbor(selected_position1)
        self.draw(generate_hint(selected_position2, 1), [selected_position1, selected_position2])
        selecting = True
        while selecting:
            key = self.stdscr.getch()

            # If key is an arrow key then move the selected position
            # If key is enter then set selected position
            match key:
                case curses.KEY_UP:
                    if self.is_point_inside_grid(selected_position1 + np.array([-1, 0])):
                        selected_position2 = selected_position1 + np.array([-1, 0])
                case curses.KEY_DOWN:
                    if self.is_point_inside_grid(selected_position1 + np.array([1, 0])):
                        selected_position2 = selected_position1 + np.array([1, 0])
                case curses.KEY_LEFT:
                    if self.is_point_inside_grid(selected_position1 + np.array([0, -1])):
                        selected_position2 = selected_position1 + np.array([0, -1])
                case curses.KEY_RIGHT:
                    if self.is_point_inside_grid(selected_position1 + np.array([0, 1])):
                        selected_position2 = selected_position1 + np.array([0, 1])
                case curses.KEY_ENTER | 10 | 13:
                    selecting = False

            self.draw(generate_hint(selected_position2, 1), [selected_position1, selected_position2])

        return selected_position1, selected_position2

    def draw(self, hint, selected_positions=None):
        self.stdscr.clear()

        for (i, j), symbol in np.ndenumerate(self.grid):
            try:
                color_pair_index = self.symbols.index(symbol) + 1
            except ValueError:
                color_pair_index = 0

            if selected_positions is not None\
                    and any(((i, j) == selected_position).all() for selected_position in selected_positions):
                color_pair_index += len(self.colors)

            self.stdscr.addstr(i, j * 2, symbol, curses.color_pair(color_pair_index))

            if j < self.grid.shape[1] - 1:
                self.stdscr.addstr(i, (j * 2) + 1, Game.empty_symbol)

        if hint is not None:
            self.stdscr.addstr(self.grid.shape[1] + 2, 0, hint)

    def refresh_grid(self):
        self.stdscr.timeout(100)
        skip_hint = "Press any key to skip..."
        self.draw(skip_hint)

        skipped = False
        while True:
            match = self.find_match()

            if match is None:
                break

            (position, direction) = match

            # self.stdscr.clear()
            # self.stdscr.addstr(12, 0, str(position))
            # self.stdscr.addstr(13, 0, str(direction))

            # TODO: Allow option of removing straight lines instead of whole cluster
            # A cluster is a group of adjacent points with the same symbol
            points_to_remove = self.find_cluster(position)

            # Set points to empty
            for point in points_to_remove:
                self.grid[point[0], point[1]] = Game.empty_symbol

            self.draw(skip_hint)
            empty_points = points_to_remove
            while empty_points.size > 0:
                self.shift(empty_points)
                empty_points = self.find_empty_points()
                self.draw(skip_hint)

                if not skipped:
                    skipped = self.stdscr.getch() != -1

            if not skipped:
                skipped = self.stdscr.getch() != -1

        # Remove timeout - make input blocking without a timeout
        self.stdscr.timeout(-1)

    def find_matches_in_direction(self, direction, match_length):
        slice_top_left_indices = np.stack([
            np.arange(match_length) if direction[0] != 0 else np.zeros(match_length, int),
            np.arange(match_length) if direction[1] != 0 else np.zeros(match_length, int)
        ])
        slice_bottom_right_indices = np.stack([
            slice_top_left_indices[0] + self.grid.shape[0] - match_length + 1
            if direction[0] != 0 else np.full(match_length, self.grid.shape[0], int),
            slice_top_left_indices[1] + self.grid.shape[1] - match_length + 1
            if direction[1] != 0 else np.full(match_length, self.grid.shape[1], int),
        ])
        slices = np.stack([self.grid[slice_top_left_indices[0, i]:slice_bottom_right_indices[0, i],
                           slice_top_left_indices[1, i]:slice_bottom_right_indices[1, i]]
                           for i in range(slice_top_left_indices.shape[1])])
        adjacent_slices_equal = np.stack([slices[i] == slices[i + 1] for i in range(slices.shape[0] - 1)])
        adjacent_slices_equal_and = np.logical_and.reduce(adjacent_slices_equal)
        x = np.argwhere(adjacent_slices_equal_and)

        # self.stdscr.clear()
        # self.stdscr.addstr(11, 0, str(x))
        # self.draw()
        # self.stdscr.getch()

        return x

    def find_match(self, match_length=3):
        directions = np.array([
            (1, 0),
            (0, 1),
            (1, 1)
        ])
        for direction in directions:
            matches = self.find_matches_in_direction(direction, match_length)
            if matches.size > 0:
                match = matches[0, :]  # Get first match position
                return match, direction
        return None

    def find_cluster(self, point):
        fringe = set()
        visited = set()
        current_point = point

        # self.stdscr.clear()
        # self.stdscr.addstr(11, 0, str(current_point))
        # self.draw()
        # self.stdscr.getch()

        neighbors = set(map(tuple, self.find_same_symbol_neighbors(current_point)))
        fringe |= neighbors
        while len(fringe) != 0:
            current_point = np.array(fringe.pop())
            visited.add(tuple(current_point))
            neighbors = set(map(tuple, self.find_same_symbol_neighbors(current_point)))
            fringe |= neighbors
            fringe -= visited

            # self.stdscr.clear()
            # self.stdscr.addstr(11, 0, str(current_point))
            # self.stdscr.addstr(12, 0, str(neighbors))
            # self.stdscr.addstr(13, 0, str(fringe))
            # self.stdscr.addstr(14, 0, str(visited))
            # self.draw()
            # self.stdscr.getch()

        return np.array(list(visited))

    def find_same_symbol_neighbors(self, point):
        # self.stdscr.clear()
        # self.stdscr.addstr(11, 0, str(point))
        same_symbol_neighbours = np.stack([
            point + translation
            # For all translations [(-1, -1), (-1, 0), (-1, 1), ..., (1, 1)]
            for translation in np.mgrid[-1:2, -1:2].transpose(1, 2, 0).reshape(-1, 2)
            # If:
            # (1) Resultant point is different from given point (i.e. translation is not [0, 0])
            # (2) Resultant point is inside the grid
            # (3) Resultant point has same symbol as given point
            if not (translation == np.zeros(2)).all() and self.is_point_inside_grid(point + translation)
            and self.grid[point[0], point[1]] == self.grid[(point + translation)[0], (point + translation)[1]]
        ])
        # self.stdscr.addstr(12, 0, str(same_symbol_neighbours))
        return same_symbol_neighbours

    def is_point_inside_grid(self, point):
        return (point >= np.zeros(2)).all() and (point < np.array(self.grid.shape)).all()

    def find_empty_points(self):
        # Return indices of elements that are the empty symbol
        return np.argwhere(self.grid == Game.empty_symbol)

    def shift(self, empty_points):
        # For each column that contains (one or more) empty points
        for j in np.unique(empty_points[:, 1]):
            # Find lowest point in that column
            rows = empty_points[empty_points[:, 1] == j, 0]  # Rows of empty points in column j
            row_max = rows.max()

            for i in range(row_max, 0, -1):
                self.grid[i, j] = self.grid[i - 1, j]

            # Fill empty point at top with a random symbol
            self.grid[0, j] = Game.get_random_symbol()

    def swap(self, selected_point1, selected_point2):
        self.grid[selected_point1[0], selected_point1[1]], self.grid[selected_point2[0], selected_point2[1]] = \
            self.grid[selected_point2[0], selected_point2[1]], self.grid[selected_point1[0], selected_point1[1]]


def main(stdscr):
    random.seed(1234)  # For debugging

    # Initialise colours
    colors = [
        curses.COLOR_WHITE,
        curses.COLOR_CYAN,
        curses.COLOR_MAGENTA,
        curses.COLOR_GREEN,
        curses.COLOR_RED,
        curses.COLOR_YELLOW
    ]
    for i, color in enumerate(colors):
        curses.init_pair(i + 1, color, curses.COLOR_BLACK)
        curses.init_pair(len(colors) + i + 1, curses.COLOR_BLACK, color)

    # Show cursor
    curses.curs_set(1)

    # TODO: Title screen

    # Hide cursor
    curses.curs_set(0)
    game = Game(stdscr, colors)
    game.play()

    # Show cursor
    curses.curs_set(1)

    # TODO: "Game over" screen


if __name__ == "__main__":
    curses.wrapper(main)

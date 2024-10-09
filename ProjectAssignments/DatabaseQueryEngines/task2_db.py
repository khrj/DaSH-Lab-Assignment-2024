'''
A simple single-table database implementation with insert, select, update, and
delete operations. This can be made more object-oriented by having a Row class,
with its own methods to update or delete its own values.

A Table class can also extend this implementation to have multiple tables in a
single Database: in fact, the current implemntation of Database is almost like
an implementation of a Table class in itself.

Adding a Table class would further allow for methods like joins. Since this is a
minimal demonstration, I haven't added more complex functions, but their
implementation, on a non-optimised level, should be straightforward.
'''
class Database:
    rows = []
    column_name_index = {}
    column_types = {}

    def __init__(self, columns) -> None:
        for i, column in enumerate(columns):
            self.column_name_index[column] = i

    # insert method with validation
    def insert(self, rows):
        if len(self.rows) == 0:
            # setup types for each column
            for i, value in enumerate(rows[0]):
                self.column_types[i] = type(value)

        for i, row in enumerate(rows):
            for j, value in enumerate(row):
                if type(value) != self.column_types[j]:
                    raise ValueError(f"Invalid type for column {j} in row {i}")
            self.rows.append(row)

    # single column select method
    # returns 2x2 matrix of rows
    def select(self, column, condition_type="=", value=None):
        column_index = self.column_name_index[column]

        if value is None:
            return self.rows
        elif condition_type == "=":
            return [row for row in self.rows if row[column_index] == value]
        elif condition_type == ">":
            return [row for row in self.rows if row[column_index] > value]
        elif condition_type == "<":
            return [row for row in self.rows if row[column_index] < value]
        elif condition_type == ">=":
            return [row for row in self.rows if row[column_index] >= value]
        elif condition_type == "<=":
            return [row for row in self.rows if row[column_index] <= value]
        elif condition_type == "!=":
            return [row for row in self.rows if row[column_index] != value]

    # update method based on a single column
    def update(
        self,
        set_column,
        set_value,
        where_column,
        where_condition_type="=",
        where_value=None,
    ):
        set_column_index = self.column_name_index[set_column]
        where_column_index = self.column_name_index[where_column]

        new_rows = []

        for row in self.rows:
            if where_value is not None:
                if (
                    (
                        where_condition_type == "="
                        and row[where_column_index] == where_value
                    )
                    or (
                        where_condition_type == "!="
                        and row[where_column_index] != where_value
                    )
                    or (
                        where_condition_type == "<"
                        and row[where_column_index] < where_value
                    )
                    or (
                        where_condition_type == ">"
                        and row[where_column_index] > where_value
                    )
                    or (
                        where_condition_type == "<="
                        and row[where_column_index] <= where_value
                    )
                    or (
                        where_condition_type == ">="
                        and row[where_column_index] >= where_value
                    )
                ):
                    row[set_column_index] = set_value

            new_rows.append(row)

        self.rows = new_rows

    # delete method based on a single column
    def delete(self, column, condition_type="=", value=None):
        condition_conjugate = {
            "=": "!=",
            "!=": "=",
            "<": ">=",
            ">": "<=",
            "<=": ">",
            ">=": "<",
        }

        self.rows = self.select(column, condition_conjugate[condition_type], value)


# Test implementation
if __name__ == "__main__":
    db = Database(["name", "age"])
    db.insert([["A", 25], ["B", 30], ["C", 25]])

    try:
        db.insert([["D", "25"]])
        assert False  # insert should generate a ValueError due to a type mismatch
    except ValueError:
        pass

    assert len(db.rows) == 3
    assert db.select("age", ">", 25) == [["B", 30]]
    assert db.select("age", "=", 25) == [["A", 25], ["C", 25]]

    db.update("age", 26, "name", "=", "A")
    assert db.select("name", "=", "A") == [["A", 26]]

    db.delete("age", ">", 25)
    assert db.rows == [["C", 25]]

    print("Test cases passed.")

import uuid
from datetime import datetime
import random
import json
import os

# TODO: Implement user login with username/password
class Task:
    def __init__(self, title, description, due_date):
        self.id = str(uuid.uuid4())
        self.title = title
        self.description = description
        self.due_date = due_date
        self.completed = False
        self.created_at = datetime.now()

    def mark_complete(self):
        self.completed = True

    def to_dict(self):
        return {
            "id": self.id,
            "title": self.title,
            "description": self.description,
            "due_date": self.due_date.isoformat(),
            "completed": self.completed,
            "created_at": self.created_at.isoformat()
        }

    @staticmethod
    def from_dict(data):
        task = Task(
            title=data["title"],
            description=data["description"],
            due_date=datetime.fromisoformat(data["due_date"])
        )
        task.id = data["id"]
        task.completed = data["completed"]
        task.created_at = datetime.fromisoformat(data["created_at"])
        return task

class TaskManager:
    def __init__(self, storage_path="tasks.json"):
        self.tasks = []
        self.storage_path = storage_path
        self.load_tasks()

    def add_task(self, title, description, due_date):
        task = Task(title, description, due_date)
        self.tasks.append(task)
        self.save_tasks()

    def remove_task(self, task_id):
        self.tasks = [task for task in self.tasks if task.id != task_id]
        self.save_tasks()

    def get_pending_tasks(self):
        return [task for task in self.tasks if not task.completed]

    def get_completed_tasks(self):
        return [task for task in self.tasks if task.completed]

    def complete_task(self, task_id):
        for task in self.tasks:
            if task.id == task_id:
                task.mark_complete()
                self.save_tasks()
                break

    def list_tasks(self):
        for task in self.tasks:
            status = "✓" if task.completed else "✗"
            print(f"[{status}] {task.title} (Due: {task.due_date.date()})")

    def save_tasks(self):
        with open(self.storage_path, "w") as f:
            json.dump([task.to_dict() for task in self.tasks], f, indent=2)

    def load_tasks(self):
        if os.path.exists(self.storage_path):
            with open(self.storage_path, "r") as f:
                data = json.load(f)
                self.tasks = [Task.from_dict(item) for item in data]

# TODO: Implement argument parsing for CLI with argparse
def generate_fake_tasks(manager, count=5):
    titles = ["Walk dog", "Do taxes", "Write blog post", "Study math", "Call parents"]
    for _ in range(count):
        title = random.choice(titles)
        description = f"{title} description"
        due_date = datetime.now()
        manager.add_task(title, description, due_date)

# TODO: Add error handling for invalid task IDs
def simulate_user_interaction():
    manager = TaskManager()
    generate_fake_tasks(manager, 3)

    print("All Tasks:")
    manager.list_tasks()

    if manager.tasks:
        manager.complete_task(manager.tasks[0].id)

    print("\nPending:")
    for task in manager.get_pending_tasks():
        print(task.title)

    print("\nCompleted:")
    for task in manager.get_completed_tasks():
        print(task.title)

# TODO: Add logging with log levels
def main():
    simulate_user_interaction()

# TODO: Refactor to support multiple users with different task files
if __name__ == "__main__":
    main()

import uuid
from datetime import datetime
import random

# TODO: Add user authentication
class Task:
    def __init__(self, title, description, due_date):
        self.id = uuid.uuid4()
        self.title = title
        self.description = description
        self.due_date = due_date
        self.completed = False
        self.created_at = datetime.now()

    def mark_complete(self):
        self.completed = True

    def __str__(self):
        return f"{self.title} - {'Done' if self.completed else 'Pending'}"

# TODO: Save tasks to file
class TaskManager:
    def __init__(self):
        self.tasks = []

    def add_task(self, title, description, due_date):
        task = Task(title, description, due_date)
        self.tasks.append(task)

    def remove_task(self, task_id):
        self.tasks = [task for task in self.tasks if task.id != task_id]

    def get_pending_tasks(self):
        return [task for task in self.tasks if not task.completed]

    def get_completed_tasks(self):
        return [task for task in self.tasks if task.completed]

    def complete_task(self, task_id):
        for task in self.tasks:
            if task.id == task_id:
                task.mark_complete()
                break

    def list_tasks(self):
        for task in self.tasks:
            print(task)

# TODO: Add CLI support
def generate_fake_tasks(manager, count=10):
    titles = ["Buy groceries", "Read a book", "Call Alice", "Fix the bike", "Water plants"]
    for _ in range(count):
        title = random.choice(titles)
        description = f"{title} description"
        due_date = datetime.now()
        manager.add_task(title, description, due_date)

def simulate_user_interaction():
    manager = TaskManager()
    generate_fake_tasks(manager, 5)
    print("All Tasks:")
    manager.list_tasks()

    print("\nCompleting first task...\n")
    if manager.tasks:
        manager.complete_task(manager.tasks[0].id)

    print("Pending Tasks:")
    for t in manager.get_pending_tasks():
        print(t)

    print("\nCompleted Tasks:")
    for t in manager.get_completed_tasks():
        print(t)

# TODO: Handle exceptions properly
def main():
    simulate_user_interaction()

# TODO: Add logging for debugging
if __name__ == "__main__":
    main()

# TODO: Implement recurring tasks
# TODO: Integrate with calendar API
# TODO: Create a GUI using Tkinter
# TODO: Add unit tests
# TODO: Use a database instead of in-memory list
# TODO: Schedule email reminders
# TODO: Implement dark mode for GUI
# TODO: Support multi-user accounts
# TODO: Encrypt saved data
# TODO: Optimize task sorting
# TODO: Support task categories and labels
# TODO: Refactor TaskManager class
# TODO: Use environment variables for config
# TODO: Improve input validation
# TODO: Add voice command interface
# TODO: Create mobile app interface
# TODO: Build a REST API for remote access
# TODO: Add support for attachments in tasks
# TODO: Track task history and edits
# TODO: Auto-save tasks at regular intervals
# TODO: Export tasks to CSV and JSON
# TODO: Add search functionality
# TODO: Generate daily task summary
# TODO: Monitor productivity stats
# TODO: Integrate Pomodoro timer
# TODO: Add undo functionality
# TODO: Track time spent on tasks
# TODO: Make UI responsive for different screens
# TODO: Translate app to multiple languages
# TODO: Add shortcut keys for power users
# TODO: Detect overdue tasks
# TODO: Create onboarding tutorial
# TODO: Add emoji support in task titles
# TODO: Provide feedback form
# TODO: Display motivational quotes
# TODO: Use color coding for priorities
# TODO: Add swipe gesture support
# TODO: Integrate with Slack and Discord
# TODO: Make app theme customizable
# TODO: Handle leap years in date logic
# TODO: Build weekly/monthly reports
# TODO: Use asyncio for background tasks
# TODO: Support markdown in task description
# TODO: Archive old completed tasks
# TODO: Implement AI-based task suggestions
# TODO: Create plugin system for extensions
# TODO: Visualize tasks on calendar view
# TODO: Set location-based reminders
# TODO: Add backup and restore functionality
# TODO: Prevent duplicate task entries
# TODO: Add keyboard navigation
# TODO: Send push notifications
# TODO: Link related tasks together

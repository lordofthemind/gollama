# **Project: Gollama CLI**

## **Purpose**
A command-line tool (`gollama`) that integrates directly with the Ollama API, allowing developers to interact with LLMs (like language models) directly from their terminal or VS Code workspace without leaving the IDE. `Gollama` will manage prompts and responses in an organized way, making it easy to access past interactions and keep code, notes, and ideas in sync within the project environment.

---

## **Features and Workflow**

1. **Initialization and Setup**
   - **Command**: `gollama init`
   - **Function**: Initializes `Gollama` in the current project workspace.
   - **Folder Structure**:
     - Creates a `.gollama` directory at the root of the project containing two folders:
       - `prompts/`: Where the user writes their prompts.
       - `deliveries/`: Where `Gollama` saves Ollama’s responses.
     - Inside `.gollama`, create a `config.json` file that holds initial configurations (e.g., default model and other settings).
   - **Configuration Prompt**:
     - Upon initialization, `Gollama` prompts the user for:
       - Default model (e.g., `GPT-3.5`, `GPT-4`, etc.).
       - Default response format (markdown or plain text).
       - Options for response length, temperature, and other model parameters if supported by the Ollama backend.

2. **Prompt Management and File-Based Interface**
   - **Command**: `gollama prompt <filename>`
   - **Function**: Opens or creates a new prompt file in the `prompts/` directory where the user can write their query.
   - **File Format**:
     - Each prompt file can be a `.txt` or `.md` file.
     - Prompts are saved with timestamps in the `prompts/` folder for easy organization and history tracking.

3. **Sending Prompts and Receiving Responses**
   - **Command**: `gollama send <filename>`
   - **Function**: Sends the content of a specified prompt file to the configured Ollama model and fetches the response.
   - **Response Handling**:
     - The response is saved in `deliveries/` as a new markdown file, named after the original prompt file but appended with a timestamp or version.
     - For example, if the prompt file is `prompt1.txt`, the response file could be saved as `prompt1_2024-11-04.md`.
   - **Automatic Opening**: If used within VS Code, responses can be automatically opened as a preview in a new editor pane, letting users view Ollama's answers in real-time without leaving the IDE.

4. **Interactive Prompting and Response Options**
   - **Command**: `gollama chat`
   - **Function**: Provides a more interactive chat-like experience where users can enter prompts directly in the terminal without saving them first. Each interaction is logged, with prompts and responses saved to `prompts/` and `deliveries/` respectively.
   - **Optional Arguments**: For one-off prompts or configurations:
     - `--model`: Override the default model for a specific query.
     - `--temperature`, `--max-tokens`, etc.: Modify model behavior for individual queries.

5. **Enhanced Configuration Management**
   - **Command**: `gollama config`
   - **Function**: Opens a configuration wizard to modify settings after initialization.
   - **Editable Options**:
     - Model settings, default response formats, and user preferences.
     - Option to set up environment-specific configurations, such as toggling between staging and production models or configuring proxies if necessary.

---

## **Additional Features**

1. **Template Prompts**
   - `gollama` could include a set of pre-defined templates for common prompts. These templates could be stored in a `templates/` folder within `.gollama` and accessible through `gollama templates list` or `gollama templates use <template-name>`. This would help users quickly set up common queries or workflow-specific tasks.

2. **Automatic Summaries and Highlights**
   - Enable an option for automatic summaries of the response, providing the main points of long outputs or adding highlights to the response files.

3. **Search and History Commands**
   - `gollama history`: Lists past prompt files and their corresponding response files, helping users quickly find previous interactions.
   - `gollama search <keyword>`: Searches through all previous prompts and responses for specific keywords or terms.

4. **Context-Aware Prompting**
   - If used in a specific project, `Gollama` could optionally pass metadata about the project’s current files or code context (with user permission). This way, prompts could be more relevant and specific to the project at hand.

5. **VS Code Integration with Custom Commands**
   - **Task Runner or Extension**: Consider creating a `Gollama` VS Code extension with custom commands like `Gollama: Send Prompt`, `Gollama: Open Deliveries`, and `Gollama: Config` for a more seamless IDE experience.
   - This extension could also display prompt history and responses in a sidebar, enabling fast navigation.

---

## **Example Workflow**

1. **Initialize the Tool**:  
   ```bash
   gollama init
   ```
   - User is prompted to choose a default model and preferences.

2. **Create a Prompt**:  
   ```bash
   gollama prompt my_prompt.md
   ```
   - User writes the prompt in `prompts/my_prompt.md`.

3. **Send Prompt and Receive Response**:  
   ```bash
   gollama send my_prompt.md
   ```
   - `Gollama` sends the prompt to Ollama, receives the response, and saves it as `deliveries/my_prompt_2024-11-04.md`.

4. **Review in VS Code**:  
   - The response file opens automatically in a new editor pane for easy reading and reference.

   
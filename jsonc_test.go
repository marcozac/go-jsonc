// Copyright 2023 Marco Zaccaro. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsonc

import (
	_ "embed"
	"testing"

	"github.com/marcozac/go-jsonc/internal/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/small.json
	_small []byte

	//go:embed testdata/small_uncommented.json
	_smallUncommented []byte

	//go:embed testdata/small_no_comment_runes.json
	_smallNoCommentRunes []byte

	//go:embed testdata/medium.json
	_medium []byte

	//go:embed testdata/medium_uncommented.json
	_mediumUncommented []byte

	//go:embed testdata/medium_no_comment_runes.json
	_mediumNoCommentRunes []byte

	_invalidChar = []byte("\xa5")
)

func FieldsValue[T DataType](t require.TestingT, j T) {
	switch j := any(j).(type) {
	case Small:
		var w Small
		require.NoError(t, json.Unmarshal(_smallUncommented, &w), "unmarshal json without comments failed")
		assert.Equal(t, w, j, "unmarshaled JSON is invalid")
		w.X = "x" // ensure fields are checked
		assert.NotEqual(t, w, j, "not all fields were checked")
	case Medium:
		var w Medium
		require.NoError(t, json.Unmarshal(_mediumUncommented, &w), "unmarshal json without comments failed")
		require.Equal(t, w, j, "unmarshaled JSON is invalid")
		w.CSS.EditorSuggestInsertMode = "insert_replace" // ensure fields are checked
		assert.NotEqual(t, w, j, "not all fields were checked")
	case SmallNoCommentRunes:
		var w SmallNoCommentRunes
		require.NoError(t, json.Unmarshal(_smallNoCommentRunes, &w), "unmarshal json without comments failed")
		assert.Equal(t, w, j, "unmarshaled JSON is invalid")
		w.X = "x" // ensure fields are checked
		assert.NotEqual(t, w, j, "not all fields were checked")
	case MediumNoCommentRunes:
		var w MediumNoCommentRunes
		require.NoError(t, json.Unmarshal(_mediumNoCommentRunes, &w), "unmarshal json without comments failed")
		require.Equal(t, w, j, "unmarshaled JSON is invalid")
		w.CSS.EditorSuggestInsertMode = "insert_replace" // ensure fields are checked
		assert.NotEqual(t, w, j, "not all fields were checked")
	default:
		assert.Fail(t, "unexpected data type: %T", j)
	}
}

func TestHasCommentRunes(t *testing.T) {
	t.Parallel()
	for _, tt := range hasCommentRunesTests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.Want, HasCommentRunes(tt.Data))
		})
	}
}

var hasCommentRunesTests = [...]struct {
	Name string
	Data []byte
	Want bool
}{
	{"Small/Commented", _small, true},
	{"Small/Uncommented", _smallUncommented, true},
	{"Small/NoCommentRunes", _smallNoCommentRunes, false},
	{"Medium/Commented", _medium, true},
	{"Medium/Uncommented", _mediumUncommented, true},
	{"Medium/NoCommentRunes", _mediumNoCommentRunes, false},
}

func BenchmarkHasCommentRunes(b *testing.B) {
	for _, tt := range hasCommentRunesTests {
		tt := tt
		b.Run(tt.Name, func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					assert.Equal(b, tt.Want, HasCommentRunes(tt.Data))
				}
			})
		})
	}
}

type DataType interface {
	Small | SmallNoCommentRunes | Medium | MediumNoCommentRunes
}

type SmallNoCommentRunes Small

type Small struct {
	Foo   string `json:"foo"`
	Baz   string `json:"baz"`
	Hello string `json:"hello"`
	X     string `json:"x,omitempty"`
}

type MediumNoCommentRunes Medium

type Medium struct {
	TelemetryTelemetryLevel                                         string   `json:"telemetry.telemetryLevel,omitempty"`
	DiffEditorCodeLens                                              bool     `json:"diffEditor.codeLens,omitempty"`
	DiffEditorDiffAlgorithm                                         string   `json:"diffEditor.diffAlgorithm,omitempty"`
	DiffEditorIgnoreTrimWhitespace                                  bool     `json:"diffEditor.ignoreTrimWhitespace,omitempty"`
	DiffEditorMaxComputationTime                                    int      `json:"diffEditor.maxComputationTime,omitempty"`
	DiffEditorRenderIndicators                                      bool     `json:"diffEditor.renderIndicators,omitempty"`
	DiffEditorRenderMarginRevertIcon                                bool     `json:"diffEditor.renderMarginRevertIcon,omitempty"`
	DiffEditorRenderSideBySide                                      bool     `json:"diffEditor.renderSideBySide,omitempty"`
	DiffEditorWordWrap                                              string   `json:"diffEditor.wordWrap,omitempty"`
	EditorAcceptSuggestionOnCommitCharacter                         bool     `json:"editor.acceptSuggestionOnCommitCharacter,omitempty"`
	EditorAcceptSuggestionOnEnter                                   string   `json:"editor.acceptSuggestionOnEnter,omitempty"`
	EditorAccessibilitySupport                                      string   `json:"editor.accessibilitySupport,omitempty"`
	EditorAutoClosingBrackets                                       string   `json:"editor.autoClosingBrackets,omitempty"`
	EditorAutoClosingDelete                                         string   `json:"editor.autoClosingDelete,omitempty"`
	EditorAutoClosingOvertype                                       string   `json:"editor.autoClosingOvertype,omitempty"`
	EditorAutoClosingQuotes                                         string   `json:"editor.autoClosingQuotes,omitempty"`
	EditorAutoIndent                                                string   `json:"editor.autoIndent,omitempty"`
	EditorAutoSurround                                              string   `json:"editor.autoSurround,omitempty"`
	EditorBracketPairColorizationEnabled                            bool     `json:"editor.bracketPairColorization.enabled,omitempty"`
	EditorBracketPairColorizationIndependentColorPoolPerBracketType bool     `json:"editor.bracketPairColorization.independentColorPoolPerBracketType,omitempty"`
	EditorCodeActionsOnSave                                         struct{} `json:"editor.codeActionsOnSave,omitempty"`
	EditorCodeActionWidgetShowHeaders                               bool     `json:"editor.codeActionWidget.showHeaders,omitempty"`
	EditorCodeLens                                                  bool     `json:"editor.codeLens,omitempty"`
	EditorCodeLensFontFamily                                        string   `json:"editor.codeLensFontFamily,omitempty"`
	EditorCodeLensFontSize                                          int      `json:"editor.codeLensFontSize,omitempty"`
	EditorColorDecorators                                           bool     `json:"editor.colorDecorators,omitempty"`
	EditorColorDecoratorsActivatedOn                                string   `json:"editor.colorDecoratorsActivatedOn,omitempty"`
	EditorColorDecoratorsLimit                                      int      `json:"editor.colorDecoratorsLimit,omitempty"`
	EditorColumnSelection                                           bool     `json:"editor.columnSelection,omitempty"`
	EditorCommentsIgnoreEmptyLines                                  bool     `json:"editor.comments.ignoreEmptyLines,omitempty"`
	EditorCommentsInsertSpace                                       bool     `json:"editor.comments.insertSpace,omitempty"`
	EditorCopyWithSyntaxHighlighting                                bool     `json:"editor.copyWithSyntaxHighlighting,omitempty"`
	EditorCursorBlinking                                            string   `json:"editor.cursorBlinking,omitempty"`
	EditorCursorSmoothCaretAnimation                                string   `json:"editor.cursorSmoothCaretAnimation,omitempty"`
	EditorCursorStyle                                               string   `json:"editor.cursorStyle,omitempty"`
	EditorCursorSurroundingLines                                    int      `json:"editor.cursorSurroundingLines,omitempty"`
	EditorCursorSurroundingLinesStyle                               string   `json:"editor.cursorSurroundingLinesStyle,omitempty"`
	EditorCursorWidth                                               int      `json:"editor.cursorWidth,omitempty"`
	EditorDefaultColorDecorators                                    bool     `json:"editor.defaultColorDecorators,omitempty"`
	EditorDefaultFoldingRangeProvider                               any      `json:"editor.defaultFoldingRangeProvider,omitempty"`
	EditorDefaultFormatter                                          any      `json:"editor.defaultFormatter,omitempty"`
	EditorDefinitionLinkOpensInPeek                                 bool     `json:"editor.definitionLinkOpensInPeek,omitempty"`
	EditorDetectIndentation                                         bool     `json:"editor.detectIndentation,omitempty"`
	EditorDragAndDrop                                               bool     `json:"editor.dragAndDrop,omitempty"`
	EditorDropIntoEditorEnabled                                     bool     `json:"editor.dropIntoEditor.enabled,omitempty"`
	EditorDropIntoEditorShowDropSelector                            string   `json:"editor.dropIntoEditor.showDropSelector,omitempty"`
	EditorEmptySelectionClipboard                                   bool     `json:"editor.emptySelectionClipboard,omitempty"`
	EditorFastScrollSensitivity                                     int      `json:"editor.fastScrollSensitivity,omitempty"`
	EditorFindAddExtraSpaceOnTop                                    bool     `json:"editor.find.addExtraSpaceOnTop,omitempty"`
	EditorFindAutoFindInSelection                                   string   `json:"editor.find.autoFindInSelection,omitempty"`
	EditorFindCursorMoveOnType                                      bool     `json:"editor.find.cursorMoveOnType,omitempty"`
	EditorFindGlobalFindClipboard                                   bool     `json:"editor.find.globalFindClipboard,omitempty"`
	EditorFindLoop                                                  bool     `json:"editor.find.loop,omitempty"`
	EditorFindSeedSearchStringFromSelection                         string   `json:"editor.find.seedSearchStringFromSelection,omitempty"`
	EditorFolding                                                   bool     `json:"editor.folding,omitempty"`
	EditorFoldingHighlight                                          bool     `json:"editor.foldingHighlight,omitempty"`
	EditorFoldingImportsByDefault                                   bool     `json:"editor.foldingImportsByDefault,omitempty"`
	EditorFoldingMaximumRegions                                     int      `json:"editor.foldingMaximumRegions,omitempty"`
	EditorFoldingStrategy                                           string   `json:"editor.foldingStrategy,omitempty"`
	EditorFontFamily                                                string   `json:"editor.fontFamily,omitempty"`
	EditorFontLigatures                                             bool     `json:"editor.fontLigatures,omitempty"`
	EditorFontSize                                                  int      `json:"editor.fontSize,omitempty"`
	EditorFontVariations                                            bool     `json:"editor.fontVariations,omitempty"`
	EditorFontWeight                                                string   `json:"editor.fontWeight,omitempty"`
	EditorFormatOnPaste                                             bool     `json:"editor.formatOnPaste,omitempty"`
	EditorFormatOnSave                                              bool     `json:"editor.formatOnSave,omitempty"`
	EditorFormatOnSaveMode                                          string   `json:"editor.formatOnSaveMode,omitempty"`
	EditorFormatOnType                                              bool     `json:"editor.formatOnType,omitempty"`
	EditorGlyphMargin                                               bool     `json:"editor.glyphMargin,omitempty"`
	EditorGotoLocationAlternativeDeclarationCommand                 string   `json:"editor.gotoLocation.alternativeDeclarationCommand,omitempty"`
	EditorGotoLocationAlternativeDefinitionCommand                  string   `json:"editor.gotoLocation.alternativeDefinitionCommand,omitempty"`
	EditorGotoLocationAlternativeImplementationCommand              string   `json:"editor.gotoLocation.alternativeImplementationCommand,omitempty"`
	EditorGotoLocationAlternativeReferenceCommand                   string   `json:"editor.gotoLocation.alternativeReferenceCommand,omitempty"`
	EditorGotoLocationAlternativeTypeDefinitionCommand              string   `json:"editor.gotoLocation.alternativeTypeDefinitionCommand,omitempty"`
	EditorGotoLocationMultipleDeclarations                          string   `json:"editor.gotoLocation.multipleDeclarations,omitempty"`
	EditorGotoLocationMultipleDefinitions                           string   `json:"editor.gotoLocation.multipleDefinitions,omitempty"`
	EditorGotoLocationMultipleImplementations                       string   `json:"editor.gotoLocation.multipleImplementations,omitempty"`
	EditorGotoLocationMultipleReferences                            string   `json:"editor.gotoLocation.multipleReferences,omitempty"`
	EditorGotoLocationMultipleTypeDefinitions                       string   `json:"editor.gotoLocation.multipleTypeDefinitions,omitempty"`
	EditorGuidesBracketPairs                                        bool     `json:"editor.guides.bracketPairs,omitempty"`
	EditorGuidesBracketPairsHorizontal                              string   `json:"editor.guides.bracketPairsHorizontal,omitempty"`
	EditorGuidesHighlightActiveBracketPair                          bool     `json:"editor.guides.highlightActiveBracketPair,omitempty"`
	EditorGuidesHighlightActiveIndentation                          bool     `json:"editor.guides.highlightActiveIndentation,omitempty"`
	EditorGuidesIndentation                                         bool     `json:"editor.guides.indentation,omitempty"`
	EditorHideCursorInOverviewRuler                                 bool     `json:"editor.hideCursorInOverviewRuler,omitempty"`
	EditorHoverAbove                                                bool     `json:"editor.hover.above,omitempty"`
	EditorHoverDelay                                                int      `json:"editor.hover.delay,omitempty"`
	EditorHoverEnabled                                              bool     `json:"editor.hover.enabled,omitempty"`
	EditorHoverSticky                                               bool     `json:"editor.hover.sticky,omitempty"`
	EditorIndentSize                                                string   `json:"editor.indentSize,omitempty"`
	EditorInlayHintsEnabled                                         string   `json:"editor.inlayHints.enabled,omitempty"`
	EditorInlayHintsFontFamily                                      string   `json:"editor.inlayHints.fontFamily,omitempty"`
	EditorInlayHintsFontSize                                        int      `json:"editor.inlayHints.fontSize,omitempty"`
	EditorInlayHintsPadding                                         bool     `json:"editor.inlayHints.padding,omitempty"`
	EditorInlineSuggestEnabled                                      bool     `json:"editor.inlineSuggest.enabled,omitempty"`
	EditorInlineSuggestShowToolbar                                  string   `json:"editor.inlineSuggest.showToolbar,omitempty"`
	EditorInlineSuggestSuppressSuggestions                          bool     `json:"editor.inlineSuggest.suppressSuggestions,omitempty"`
	EditorInsertSpaces                                              bool     `json:"editor.insertSpaces,omitempty"`
	EditorLanguageBrackets                                          any      `json:"editor.language.brackets,omitempty"`
	EditorLanguageColorizedBracketPairs                             any      `json:"editor.language.colorizedBracketPairs,omitempty"`
	EditorLetterSpacing                                             int      `json:"editor.letterSpacing,omitempty"`
	EditorLightbulbEnabled                                          bool     `json:"editor.lightbulb.enabled,omitempty"`
	EditorLineHeight                                                int      `json:"editor.lineHeight,omitempty"`
	EditorLineNumbers                                               string   `json:"editor.lineNumbers,omitempty"`
	EditorLinkedEditing                                             bool     `json:"editor.linkedEditing,omitempty"`
	EditorLinks                                                     bool     `json:"editor.links,omitempty"`
	EditorMatchBrackets                                             string   `json:"editor.matchBrackets,omitempty"`
	EditorMinimapAutohide                                           bool     `json:"editor.minimap.autohide,omitempty"`
	EditorMinimapEnabled                                            bool     `json:"editor.minimap.enabled,omitempty"`
	EditorMinimapMaxColumn                                          int      `json:"editor.minimap.maxColumn,omitempty"`
	EditorMinimapRenderCharacters                                   bool     `json:"editor.minimap.renderCharacters,omitempty"`
	EditorMinimapScale                                              int      `json:"editor.minimap.scale,omitempty"`
	EditorMinimapShowSlider                                         string   `json:"editor.minimap.showSlider,omitempty"`
	EditorMinimapSide                                               string   `json:"editor.minimap.side,omitempty"`
	EditorMinimapSize                                               string   `json:"editor.minimap.size,omitempty"`
	EditorMouseWheelScrollSensitivity                               int      `json:"editor.mouseWheelScrollSensitivity,omitempty"`
	EditorMouseWheelZoom                                            bool     `json:"editor.mouseWheelZoom,omitempty"`
	EditorMultiCursorModifier                                       string   `json:"editor.multiCursorModifier,omitempty"`
	EditorMultiCursorPaste                                          string   `json:"editor.multiCursorPaste,omitempty"`
	EditorOccurrencesHighlight                                      bool     `json:"editor.occurrencesHighlight,omitempty"`
	EditorOverviewRulerBorder                                       bool     `json:"editor.overviewRulerBorder,omitempty"`
	EditorPaddingBottom                                             int      `json:"editor.padding.bottom,omitempty"`
	EditorPaddingTop                                                int      `json:"editor.padding.top,omitempty"`
	EditorParameterHintsCycle                                       bool     `json:"editor.parameterHints.cycle,omitempty"`
	EditorParameterHintsEnabled                                     bool     `json:"editor.parameterHints.enabled,omitempty"`
	EditorPasteAsEnabled                                            bool     `json:"editor.pasteAs.enabled,omitempty"`
	EditorPasteAsShowPasteSelector                                  string   `json:"editor.pasteAs.showPasteSelector,omitempty"`
	EditorPeekWidgetDefaultFocus                                    string   `json:"editor.peekWidgetDefaultFocus,omitempty"`
	EditorQuickSuggestions                                          struct {
		Other    string `json:"other,omitempty"`
		Comments string `json:"comments,omitempty"`
		Strings  string `json:"strings,omitempty"`
	} `json:"editor.quickSuggestions,omitempty"`
	EditorQuickSuggestionsDelay                         int      `json:"editor.quickSuggestionsDelay,omitempty"`
	EditorRenameEnablePreview                           bool     `json:"editor.rename.enablePreview,omitempty"`
	EditorRenderControlCharacters                       bool     `json:"editor.renderControlCharacters,omitempty"`
	EditorRenderFinalNewline                            string   `json:"editor.renderFinalNewline,omitempty"`
	EditorRenderLineHighlight                           string   `json:"editor.renderLineHighlight,omitempty"`
	EditorRenderLineHighlightOnlyWhenFocus              bool     `json:"editor.renderLineHighlightOnlyWhenFocus,omitempty"`
	EditorRenderWhitespace                              string   `json:"editor.renderWhitespace,omitempty"`
	EditorRoundedSelection                              bool     `json:"editor.roundedSelection,omitempty"`
	EditorRulers                                        []any    `json:"editor.rulers,omitempty"`
	EditorScreenReaderAnnounceInlineSuggestion          bool     `json:"editor.screenReaderAnnounceInlineSuggestion,omitempty"`
	EditorScrollbarHorizontal                           string   `json:"editor.scrollbar.horizontal,omitempty"`
	EditorScrollbarHorizontalScrollbarSize              int      `json:"editor.scrollbar.horizontalScrollbarSize,omitempty"`
	EditorScrollbarScrollByPage                         bool     `json:"editor.scrollbar.scrollByPage,omitempty"`
	EditorScrollbarVertical                             string   `json:"editor.scrollbar.vertical,omitempty"`
	EditorScrollbarVerticalScrollbarSize                int      `json:"editor.scrollbar.verticalScrollbarSize,omitempty"`
	EditorScrollBeyondLastColumn                        int      `json:"editor.scrollBeyondLastColumn,omitempty"`
	EditorScrollBeyondLastLine                          bool     `json:"editor.scrollBeyondLastLine,omitempty"`
	EditorSelectionClipboard                            bool     `json:"editor.selectionClipboard,omitempty"`
	EditorScrollPredominantAxis                         bool     `json:"editor.scrollPredominantAxis,omitempty"`
	EditorSelectionHighlight                            bool     `json:"editor.selectionHighlight,omitempty"`
	EditorSemanticHighlightingEnabled                   string   `json:"editor.semanticHighlighting.enabled,omitempty"`
	EditorSemanticTokenColorCustomizations              struct{} `json:"editor.semanticTokenColorCustomizations,omitempty"`
	EditorShowDeprecated                                bool     `json:"editor.showDeprecated,omitempty"`
	EditorShowFoldingControls                           string   `json:"editor.showFoldingControls,omitempty"`
	EditorShowUnused                                    bool     `json:"editor.showUnused,omitempty"`
	EditorSmartSelectSelectLeadingAndTrailingWhitespace bool     `json:"editor.smartSelect.selectLeadingAndTrailingWhitespace,omitempty"`
	EditorSmartSelectSelectSubwords                     bool     `json:"editor.smartSelect.selectSubwords,omitempty"`
	EditorSmoothScrolling                               bool     `json:"editor.smoothScrolling,omitempty"`
	EditorSnippetsCodeActionsEnabled                    bool     `json:"editor.snippets.codeActions.enabled,omitempty"`
	EditorSnippetSuggestions                            string   `json:"editor.snippetSuggestions,omitempty"`
	EditorStablePeek                                    bool     `json:"editor.stablePeek,omitempty"`
	EditorStickyScrollDefaultModel                      string   `json:"editor.stickyScroll.defaultModel,omitempty"`
	EditorStickyScrollEnabled                           bool     `json:"editor.stickyScroll.enabled,omitempty"`
	EditorStickyScrollMaxLineCount                      int      `json:"editor.stickyScroll.maxLineCount,omitempty"`
	EditorStickyTabStops                                bool     `json:"editor.stickyTabStops,omitempty"`
	EditorSuggestFilterGraceful                         bool     `json:"editor.suggest.filterGraceful,omitempty"`
	EditorSuggestInsertMode                             string   `json:"editor.suggest.insertMode,omitempty"`
	EditorSuggestLocalityBonus                          bool     `json:"editor.suggest.localityBonus,omitempty"`
	EditorSuggestMatchOnWordStartOnly                   bool     `json:"editor.suggest.matchOnWordStartOnly,omitempty"`
	EditorSuggestPreview                                bool     `json:"editor.suggest.preview,omitempty"`
	EditorSuggestSelectionMode                          string   `json:"editor.suggest.selectionMode,omitempty"`
	EditorSuggestShareSuggestSelections                 bool     `json:"editor.suggest.shareSuggestSelections,omitempty"`
	EditorSuggestShowClasses                            bool     `json:"editor.suggest.showClasses,omitempty"`
	EditorSuggestShowColors                             bool     `json:"editor.suggest.showColors,omitempty"`
	EditorSuggestShowConstants                          bool     `json:"editor.suggest.showConstants,omitempty"`
	EditorSuggestShowConstructors                       bool     `json:"editor.suggest.showConstructors,omitempty"`
	EditorSuggestShowCustomcolors                       bool     `json:"editor.suggest.showCustomcolors,omitempty"`
	EditorSuggestShowDeprecated                         bool     `json:"editor.suggest.showDeprecated,omitempty"`
	EditorSuggestShowEnumMembers                        bool     `json:"editor.suggest.showEnumMembers,omitempty"`
	EditorSuggestShowEnums                              bool     `json:"editor.suggest.showEnums,omitempty"`
	EditorSuggestShowEvents                             bool     `json:"editor.suggest.showEvents,omitempty"`
	EditorSuggestShowFields                             bool     `json:"editor.suggest.showFields,omitempty"`
	EditorSuggestShowFiles                              bool     `json:"editor.suggest.showFiles,omitempty"`
	EditorSuggestShowFolders                            bool     `json:"editor.suggest.showFolders,omitempty"`
	EditorSuggestShowFunctions                          bool     `json:"editor.suggest.showFunctions,omitempty"`
	EditorSuggestShowIcons                              bool     `json:"editor.suggest.showIcons,omitempty"`
	EditorSuggestShowInlineDetails                      bool     `json:"editor.suggest.showInlineDetails,omitempty"`
	EditorSuggestShowInterfaces                         bool     `json:"editor.suggest.showInterfaces,omitempty"`
	EditorSuggestShowIssues                             bool     `json:"editor.suggest.showIssues,omitempty"`
	EditorSuggestShowKeywords                           bool     `json:"editor.suggest.showKeywords,omitempty"`
	EditorSuggestShowMethods                            bool     `json:"editor.suggest.showMethods,omitempty"`
	EditorSuggestShowModules                            bool     `json:"editor.suggest.showModules,omitempty"`
	EditorSuggestShowOperators                          bool     `json:"editor.suggest.showOperators,omitempty"`
	EditorSuggestShowProperties                         bool     `json:"editor.suggest.showProperties,omitempty"`
	EditorSuggestShowReferences                         bool     `json:"editor.suggest.showReferences,omitempty"`
	EditorSuggestShowSnippets                           bool     `json:"editor.suggest.showSnippets,omitempty"`
	EditorSuggestShowStatusBar                          bool     `json:"editor.suggest.showStatusBar,omitempty"`
	EditorSuggestShowStructs                            bool     `json:"editor.suggest.showStructs,omitempty"`
	EditorSuggestShowTypeParameters                     bool     `json:"editor.suggest.showTypeParameters,omitempty"`
	EditorSuggestShowUnits                              bool     `json:"editor.suggest.showUnits,omitempty"`
	EditorSuggestShowUsers                              bool     `json:"editor.suggest.showUsers,omitempty"`
	EditorSuggestShowValues                             bool     `json:"editor.suggest.showValues,omitempty"`
	EditorSuggestShowVariables                          bool     `json:"editor.suggest.showVariables,omitempty"`
	EditorSuggestShowWords                              bool     `json:"editor.suggest.showWords,omitempty"`
	EditorSuggestSnippetsPreventQuickSuggestions        bool     `json:"editor.suggest.snippetsPreventQuickSuggestions,omitempty"`
	EditorSuggestFontSize                               int      `json:"editor.suggestFontSize,omitempty"`
	EditorSuggestLineHeight                             int      `json:"editor.suggestLineHeight,omitempty"`
	EditorSuggestOnTriggerCharacters                    bool     `json:"editor.suggestOnTriggerCharacters,omitempty"`
	EditorSuggestSelection                              string   `json:"editor.suggestSelection,omitempty"`
	EditorTabCompletion                                 string   `json:"editor.tabCompletion,omitempty"`
	EditorTabFocusMode                                  bool     `json:"editor.tabFocusMode,omitempty"`
	EditorTabSize                                       int      `json:"editor.tabSize,omitempty"`
	EditorTokenColorCustomizations                      struct{} `json:"editor.tokenColorCustomizations,omitempty"`
	EditorTrimAutoWhitespace                            bool     `json:"editor.trimAutoWhitespace,omitempty"`
	EditorUnfoldOnClickAfterEndOfLine                   bool     `json:"editor.unfoldOnClickAfterEndOfLine,omitempty"`
	EditorUnicodeHighlightAllowedCharacters             struct{} `json:"editor.unicodeHighlight.allowedCharacters,omitempty"`
	EditorUnicodeHighlightAllowedLocales                struct {
		Os     bool `json:"_os,omitempty"`
		Vscode bool `json:"_vscode,omitempty"`
	} `json:"editor.unicodeHighlight.allowedLocales,omitempty"`
	EditorUnicodeHighlightAmbiguousCharacters bool   `json:"editor.unicodeHighlight.ambiguousCharacters,omitempty"`
	EditorUnicodeHighlightIncludeComments     string `json:"editor.unicodeHighlight.includeComments,omitempty"`
	EditorUnicodeHighlightIncludeStrings      bool   `json:"editor.unicodeHighlight.includeStrings,omitempty"`
	EditorUnicodeHighlightInvisibleCharacters bool   `json:"editor.unicodeHighlight.invisibleCharacters,omitempty"`
	EditorUnicodeHighlightNonBasicASCII       string `json:"editor.unicodeHighlight.nonBasicASCII,omitempty"`
	EditorUnusualLineTerminators              string `json:"editor.unusualLineTerminators,omitempty"`
	EditorUseTabStops                         bool   `json:"editor.useTabStops,omitempty"`
	EditorWordBasedSuggestions                bool   `json:"editor.wordBasedSuggestions,omitempty"`
	EditorWordBasedSuggestionsMode            string `json:"editor.wordBasedSuggestionsMode,omitempty"`
	EditorWordBreak                           string `json:"editor.wordBreak,omitempty"`
	EditorWordSeparators                      string `json:"editor.wordSeparators,omitempty"`
	EditorWordWrap                            string `json:"editor.wordWrap,omitempty"`
	EditorWordWrapColumn                      int    `json:"editor.wordWrapColumn,omitempty"`
	EditorWrappingIndent                      string `json:"editor.wrappingIndent,omitempty"`
	EditorWrappingStrategy                    string `json:"editor.wrappingStrategy,omitempty"`
	ScmAlwaysShowActions                      bool   `json:"scm.alwaysShowActions,omitempty"`
	ScmAlwaysShowRepositories                 bool   `json:"scm.alwaysShowRepositories,omitempty"`
	ScmAutoReveal                             bool   `json:"scm.autoReveal,omitempty"`
	ScmCountBadge                             string `json:"scm.countBadge,omitempty"`
	ScmDefaultViewMode                        string `json:"scm.defaultViewMode,omitempty"`
	ScmDefaultViewSortKey                     string `json:"scm.defaultViewSortKey,omitempty"`
	ScmDiffDecorations                        string `json:"scm.diffDecorations,omitempty"`
	ScmDiffDecorationsGutterAction            string `json:"scm.diffDecorationsGutterAction,omitempty"`
	ScmDiffDecorationsGutterPattern           struct {
		Added    bool `json:"added,omitempty"`
		Modified bool `json:"modified,omitempty"`
	} `json:"scm.diffDecorationsGutterPattern,omitempty"`
	ScmDiffDecorationsGutterVisibility     string   `json:"scm.diffDecorationsGutterVisibility,omitempty"`
	ScmDiffDecorationsGutterWidth          int      `json:"scm.diffDecorationsGutterWidth,omitempty"`
	ScmDiffDecorationsIgnoreTrimWhitespace string   `json:"scm.diffDecorationsIgnoreTrimWhitespace,omitempty"`
	ScmInputFontFamily                     string   `json:"scm.inputFontFamily,omitempty"`
	ScmInputFontSize                       int      `json:"scm.inputFontSize,omitempty"`
	ScmProviderCountBadge                  string   `json:"scm.providerCountBadge,omitempty"`
	ScmRepositoriesSortOrder               string   `json:"scm.repositories.sortOrder,omitempty"`
	ScmRepositoriesVisible                 int      `json:"scm.repositories.visible,omitempty"`
	ScmShowActionButton                    bool     `json:"scm.showActionButton,omitempty"`
	SecurityAllowedUNCHosts                []any    `json:"security.allowedUNCHosts,omitempty"`
	SecurityRestrictUNCAccess              bool     `json:"security.restrictUNCAccess,omitempty"`
	SecurityWorkspaceTrustBanner           string   `json:"security.workspace.trust.banner,omitempty"`
	SecurityWorkspaceTrustEmptyWindow      bool     `json:"security.workspace.trust.emptyWindow,omitempty"`
	SecurityWorkspaceTrustEnabled          bool     `json:"security.workspace.trust.enabled,omitempty"`
	SecurityWorkspaceTrustStartupPrompt    string   `json:"security.workspace.trust.startupPrompt,omitempty"`
	SecurityWorkspaceTrustUntrustedFiles   string   `json:"security.workspace.trust.untrustedFiles,omitempty"`
	WorkbenchActivityBarIconClickBehavior  string   `json:"workbench.activityBar.iconClickBehavior,omitempty"`
	WorkbenchActivityBarVisible            bool     `json:"workbench.activityBar.visible,omitempty"`
	WorkbenchCloudChangesAutoResume        string   `json:"workbench.cloudChanges.autoResume,omitempty"`
	WorkbenchCloudChangesContinueOn        string   `json:"workbench.cloudChanges.continueOn,omitempty"`
	WorkbenchColorCustomizations           struct{} `json:"workbench.colorCustomizations,omitempty"`
	WorkbenchColorTheme                    string   `json:"workbench.colorTheme,omitempty"`
	WorkbenchCommandPaletteHistory         int      `json:"workbench.commandPalette.history,omitempty"`
	WorkbenchCommandPalettePreserveInput   bool     `json:"workbench.commandPalette.preserveInput,omitempty"`
	WorkbenchEditorAutoLockGroups          struct {
		Default                                bool `json:"default,omitempty"`
		WorkbenchEditorinputsSearchEditorInput bool `json:"workbench.editorinputs.searchEditorInput,omitempty"`
		WorkbenchEditorChatSession             bool `json:"workbench.editor.chatSession,omitempty"`
		JupyterNotebook                        bool `json:"jupyter-notebook,omitempty"`
		ImagePreviewPreviewEditor              bool `json:"imagePreview.previewEditor,omitempty"`
		VscodeAudioPreview                     bool `json:"vscode.audioPreview,omitempty"`
		VscodeVideoPreview                     bool `json:"vscode.videoPreview,omitempty"`
		JsProfileVisualizerCpuprofileTable     bool `json:"jsProfileVisualizer.cpuprofile.table,omitempty"`
		JsProfileVisualizerHeapprofileTable    bool `json:"jsProfileVisualizer.heapprofile.table,omitempty"`
		WorkbenchEditorsGettingStartedInput    bool `json:"workbench.editors.gettingStartedInput,omitempty"`
		TerminalEditor                         bool `json:"terminalEditor,omitempty"`
		WorkbenchInputInteractive              bool `json:"workbench.input.interactive,omitempty"`
		MainThreadWebviewMarkdownPreview       bool `json:"mainThreadWebview-markdown.preview,omitempty"`
	} `json:"workbench.editor.autoLockGroups,omitempty"`
	WorkbenchEditorCenteredLayoutAutoResize               bool   `json:"workbench.editor.centeredLayoutAutoResize,omitempty"`
	WorkbenchEditorCenteredLayoutFixedWidth               bool   `json:"workbench.editor.centeredLayoutFixedWidth,omitempty"`
	WorkbenchEditorCloseEmptyGroups                       bool   `json:"workbench.editor.closeEmptyGroups,omitempty"`
	WorkbenchEditorCloseOnFileDelete                      bool   `json:"workbench.editor.closeOnFileDelete,omitempty"`
	WorkbenchEditorDecorationsBadges                      bool   `json:"workbench.editor.decorations.badges,omitempty"`
	WorkbenchEditorDecorationsColors                      bool   `json:"workbench.editor.decorations.colors,omitempty"`
	WorkbenchEditorDefaultBinaryEditor                    string `json:"workbench.editor.defaultBinaryEditor,omitempty"`
	WorkbenchEditorDoubleClickTabToToggleEditorGroupSizes bool   `json:"workbench.editor.doubleClickTabToToggleEditorGroupSizes,omitempty"`
	WorkbenchEditorEnablePreview                          bool   `json:"workbench.editor.enablePreview,omitempty"`
	WorkbenchEditorEnablePreviewFromCodeNavigation        bool   `json:"workbench.editor.enablePreviewFromCodeNavigation,omitempty"`
	WorkbenchEditorEnablePreviewFromQuickOpen             bool   `json:"workbench.editor.enablePreviewFromQuickOpen,omitempty"`
	WorkbenchEditorFocusRecentEditorAfterClose            bool   `json:"workbench.editor.focusRecentEditorAfterClose,omitempty"`
	WorkbenchEditorHighlightModifiedTabs                  bool   `json:"workbench.editor.highlightModifiedTabs,omitempty"`
	WorkbenchEditorHistoryBasedLanguageDetection          bool   `json:"workbench.editor.historyBasedLanguageDetection,omitempty"`
	WorkbenchEditorLabelFormat                            string `json:"workbench.editor.labelFormat,omitempty"`
	WorkbenchEditorLanguageDetection                      bool   `json:"workbench.editor.languageDetection,omitempty"`
	WorkbenchEditorLanguageDetectionHints                 struct {
		UntitledEditors bool `json:"untitledEditors,omitempty"`
		NotebookEditors bool `json:"notebookEditors,omitempty"`
	} `json:"workbench.editor.languageDetectionHints,omitempty"`
	WorkbenchEditorLimitEnabled                        bool     `json:"workbench.editor.limit.enabled,omitempty"`
	WorkbenchEditorLimitExcludeDirty                   bool     `json:"workbench.editor.limit.excludeDirty,omitempty"`
	WorkbenchEditorLimitPerEditorGroup                 bool     `json:"workbench.editor.limit.perEditorGroup,omitempty"`
	WorkbenchEditorLimitValue                          int      `json:"workbench.editor.limit.value,omitempty"`
	WorkbenchEditorMouseBackForwardToNavigate          bool     `json:"workbench.editor.mouseBackForwardToNavigate,omitempty"`
	WorkbenchEditorNavigationScope                     string   `json:"workbench.editor.navigationScope,omitempty"`
	WorkbenchEditorOpenPositioning                     string   `json:"workbench.editor.openPositioning,omitempty"`
	WorkbenchEditorOpenSideBySideDirection             string   `json:"workbench.editor.openSideBySideDirection,omitempty"`
	WorkbenchEditorPinnedTabSizing                     string   `json:"workbench.editor.pinnedTabSizing,omitempty"`
	WorkbenchEditorPreferHistoryBasedLanguageDetection bool     `json:"workbench.editor.preferHistoryBasedLanguageDetection,omitempty"`
	WorkbenchEditorRestoreViewState                    bool     `json:"workbench.editor.restoreViewState,omitempty"`
	WorkbenchEditorRevealIfOpen                        bool     `json:"workbench.editor.revealIfOpen,omitempty"`
	WorkbenchEditorScrollToSwitchTabs                  bool     `json:"workbench.editor.scrollToSwitchTabs,omitempty"`
	WorkbenchEditorSharedViewState                     bool     `json:"workbench.editor.sharedViewState,omitempty"`
	WorkbenchEditorShowIcons                           bool     `json:"workbench.editor.showIcons,omitempty"`
	WorkbenchEditorShowTabs                            bool     `json:"workbench.editor.showTabs,omitempty"`
	WorkbenchEditorSplitInGroupLayout                  string   `json:"workbench.editor.splitInGroupLayout,omitempty"`
	WorkbenchEditorSplitOnDragAndDrop                  bool     `json:"workbench.editor.splitOnDragAndDrop,omitempty"`
	WorkbenchEditorSplitSizing                         string   `json:"workbench.editor.splitSizing,omitempty"`
	WorkbenchEditorTabCloseButton                      string   `json:"workbench.editor.tabCloseButton,omitempty"`
	WorkbenchEditorTabSizing                           string   `json:"workbench.editor.tabSizing,omitempty"`
	WorkbenchEditorTabSizingFixedMaxWidth              int      `json:"workbench.editor.tabSizingFixedMaxWidth,omitempty"`
	WorkbenchEditorTabSizingFixedMinWidth              int      `json:"workbench.editor.tabSizingFixedMinWidth,omitempty"`
	WorkbenchEditorTitleScrollbarSizing                string   `json:"workbench.editor.titleScrollbarSizing,omitempty"`
	WorkbenchEditorUntitledHint                        string   `json:"workbench.editor.untitled.hint,omitempty"`
	WorkbenchEditorUntitledLabelFormat                 string   `json:"workbench.editor.untitled.labelFormat,omitempty"`
	WorkbenchEditorWrapTabs                            bool     `json:"workbench.editor.wrapTabs,omitempty"`
	WorkbenchEditorAssociations                        struct{} `json:"workbench.editorAssociations,omitempty"`
	WorkbenchEditorLargeFileConfirmation               int      `json:"workbench.editorLargeFileConfirmation,omitempty"`
	WorkbenchExternalURIOpeners                        struct{} `json:"workbench.externalUriOpeners,omitempty"`
	WorkbenchFontAliasing                              string   `json:"workbench.fontAliasing,omitempty"`
	WorkbenchHoverDelay                                int      `json:"workbench.hover.delay,omitempty"`
	WorkbenchIconTheme                                 string   `json:"workbench.iconTheme,omitempty"`
	WorkbenchLayoutControlEnabled                      bool     `json:"workbench.layoutControl.enabled,omitempty"`
	WorkbenchLayoutControlType                         string   `json:"workbench.layoutControl.type,omitempty"`
	WorkbenchListDefaultFindMatchType                  string   `json:"workbench.list.defaultFindMatchType,omitempty"`
	WorkbenchListDefaultFindMode                       string   `json:"workbench.list.defaultFindMode,omitempty"`
	WorkbenchListFastScrollSensitivity                 int      `json:"workbench.list.fastScrollSensitivity,omitempty"`
	WorkbenchListHorizontalScrolling                   bool     `json:"workbench.list.horizontalScrolling,omitempty"`
	WorkbenchListMouseWheelScrollSensitivity           int      `json:"workbench.list.mouseWheelScrollSensitivity,omitempty"`
	WorkbenchListMultiSelectModifier                   string   `json:"workbench.list.multiSelectModifier,omitempty"`
	WorkbenchListOpenMode                              string   `json:"workbench.list.openMode,omitempty"`
	WorkbenchListScrollByPage                          bool     `json:"workbench.list.scrollByPage,omitempty"`
	WorkbenchListSmoothScrolling                       bool     `json:"workbench.list.smoothScrolling,omitempty"`
	WorkbenchListTypeNavigationMode                    string   `json:"workbench.list.typeNavigationMode,omitempty"`
	WorkbenchLocalHistoryEnabled                       bool     `json:"workbench.localHistory.enabled,omitempty"`
	WorkbenchLocalHistoryExclude                       struct{} `json:"workbench.localHistory.exclude,omitempty"`
	WorkbenchLocalHistoryMaxFileEntries                int      `json:"workbench.localHistory.maxFileEntries,omitempty"`
	WorkbenchLocalHistoryMaxFileSize                   int      `json:"workbench.localHistory.maxFileSize,omitempty"`
	WorkbenchLocalHistoryMergeWindow                   int      `json:"workbench.localHistory.mergeWindow,omitempty"`
	WorkbenchPanelDefaultLocation                      string   `json:"workbench.panel.defaultLocation,omitempty"`
	WorkbenchPanelOpensMaximized                       string   `json:"workbench.panel.opensMaximized,omitempty"`
	WorkbenchPreferredDarkColorTheme                   string   `json:"workbench.preferredDarkColorTheme,omitempty"`
	WorkbenchPreferredHighContrastColorTheme           string   `json:"workbench.preferredHighContrastColorTheme,omitempty"`
	WorkbenchPreferredHighContrastLightColorTheme      string   `json:"workbench.preferredHighContrastLightColorTheme,omitempty"`
	WorkbenchPreferredLightColorTheme                  string   `json:"workbench.preferredLightColorTheme,omitempty"`
	WorkbenchProductIconTheme                          string   `json:"workbench.productIconTheme,omitempty"`
	WorkbenchQuickOpenCloseOnFocusLost                 bool     `json:"workbench.quickOpen.closeOnFocusLost,omitempty"`
	WorkbenchQuickOpenPreserveInput                    bool     `json:"workbench.quickOpen.preserveInput,omitempty"`
	WorkbenchReduceMotion                              string   `json:"workbench.reduceMotion,omitempty"`
	WorkbenchSashHoverDelay                            int      `json:"workbench.sash.hoverDelay,omitempty"`
	WorkbenchSashSize                                  int      `json:"workbench.sash.size,omitempty"`
	WorkbenchSettingsApplyToAllProfiles                []any    `json:"workbench.settings.applyToAllProfiles,omitempty"`
	WorkbenchSettingsEditor                            string   `json:"workbench.settings.editor,omitempty"`
	WorkbenchSettingsEnableNaturalLanguageSearch       bool     `json:"workbench.settings.enableNaturalLanguageSearch,omitempty"`
	WorkbenchSettingsOpenDefaultKeybindings            bool     `json:"workbench.settings.openDefaultKeybindings,omitempty"`
	WorkbenchSettingsOpenDefaultSettings               bool     `json:"workbench.settings.openDefaultSettings,omitempty"`
	WorkbenchSettingsSettingsSearchTocBehavior         string   `json:"workbench.settings.settingsSearchTocBehavior,omitempty"`
	WorkbenchSettingsUseSplitJSON                      bool     `json:"workbench.settings.useSplitJSON,omitempty"`
	WorkbenchSideBarLocation                           string   `json:"workbench.sideBar.location,omitempty"`
	WorkbenchStartupEditor                             string   `json:"workbench.startupEditor,omitempty"`
	WorkbenchStatusBarVisible                          bool     `json:"workbench.statusBar.visible,omitempty"`
	WorkbenchTipsEnabled                               bool     `json:"workbench.tips.enabled,omitempty"`
	WorkbenchTreeExpandMode                            string   `json:"workbench.tree.expandMode,omitempty"`
	WorkbenchTreeIndent                                int      `json:"workbench.tree.indent,omitempty"`
	WorkbenchTreeRenderIndentGuides                    string   `json:"workbench.tree.renderIndentGuides,omitempty"`
	WorkbenchTrustedDomainsPromptInTrustedWorkspace    bool     `json:"workbench.trustedDomains.promptInTrustedWorkspace,omitempty"`
	WorkbenchViewAlwaysShowHeaderActions               bool     `json:"workbench.view.alwaysShowHeaderActions,omitempty"`
	WorkbenchWelcomePageWalkthroughsOpenOnInstall      bool     `json:"workbench.welcomePage.walkthroughs.openOnInstall,omitempty"`
	WindowAutoDetectColorScheme                        bool     `json:"window.autoDetectColorScheme,omitempty"`
	WindowAutoDetectHighContrast                       bool     `json:"window.autoDetectHighContrast,omitempty"`
	WindowClickThroughInactive                         bool     `json:"window.clickThroughInactive,omitempty"`
	WindowCloseWhenEmpty                               bool     `json:"window.closeWhenEmpty,omitempty"`
	WindowCommandCenter                                bool     `json:"window.commandCenter,omitempty"`
	WindowConfirmBeforeClose                           string   `json:"window.confirmBeforeClose,omitempty"`
	WindowCustomMenuBarAltFocus                        bool     `json:"window.customMenuBarAltFocus,omitempty"`
	WindowDialogStyle                                  string   `json:"window.dialogStyle,omitempty"`
	WindowDoubleClickIconToClose                       bool     `json:"window.doubleClickIconToClose,omitempty"`
	WindowNativeFullScreen                             bool     `json:"window.nativeFullScreen,omitempty"`
	WindowNativeTabs                                   bool     `json:"window.nativeTabs,omitempty"`
	WindowEnableMenuBarMnemonics                       bool     `json:"window.enableMenuBarMnemonics,omitempty"`
	WindowMenuBarVisibility                            string   `json:"window.menuBarVisibility,omitempty"`
	WindowNewWindowDimensions                          string   `json:"window.newWindowDimensions,omitempty"`
	WindowOpenFilesInNewWindow                         string   `json:"window.openFilesInNewWindow,omitempty"`
	WindowOpenFoldersInNewWindow                       string   `json:"window.openFoldersInNewWindow,omitempty"`
	WindowOpenWithoutArgumentsInNewWindow              string   `json:"window.openWithoutArgumentsInNewWindow,omitempty"`
	WindowRestoreFullscreen                            bool     `json:"window.restoreFullscreen,omitempty"`
	WindowRestoreWindows                               string   `json:"window.restoreWindows,omitempty"`
	WindowTitle                                        string   `json:"window.title,omitempty"`
	WindowTitleBarStyle                                string   `json:"window.titleBarStyle,omitempty"`
	WindowTitleSeparator                               string   `json:"window.titleSeparator,omitempty"`
	WindowZoomLevel                                    int      `json:"window.zoomLevel,omitempty"`
	FilesAssociations                                  struct{} `json:"files.associations,omitempty"`
	FilesAutoGuessEncoding                             bool     `json:"files.autoGuessEncoding,omitempty"`
	FilesAutoSave                                      string   `json:"files.autoSave,omitempty"`
	FilesAutoSaveDelay                                 int      `json:"files.autoSaveDelay,omitempty"`
	FilesDefaultLanguage                               string   `json:"files.defaultLanguage,omitempty"`
	FilesDialogDefaultPath                             string   `json:"files.dialog.defaultPath,omitempty"`
	FilesEnableTrash                                   bool     `json:"files.enableTrash,omitempty"`
	FilesEncoding                                      string   `json:"files.encoding,omitempty"`
	FilesEol                                           string   `json:"files.eol,omitempty"`
	FilesExclude                                       struct {
		Git      bool `json:"**/.git,omitempty"`
		Svn      bool `json:"**/.svn,omitempty"`
		Hg       bool `json:"**/.hg,omitempty"`
		CVS      bool `json:"**/CVS,omitempty"`
		DSStore  bool `json:"**/.DS_Store,omitempty"`
		ThumbsDb bool `json:"**/Thumbs.db,omitempty"`
	} `json:"files.exclude,omitempty"`
	FilesHotExit                 string   `json:"files.hotExit,omitempty"`
	FilesInsertFinalNewline      bool     `json:"files.insertFinalNewline,omitempty"`
	FilesParticipantsTimeout     int      `json:"files.participants.timeout,omitempty"`
	FilesReadonlyExclude         struct{} `json:"files.readonlyExclude,omitempty"`
	FilesReadonlyFromPermissions bool     `json:"files.readonlyFromPermissions,omitempty"`
	FilesReadonlyInclude         struct{} `json:"files.readonlyInclude,omitempty"`
	FilesRefactoringAutoSave     bool     `json:"files.refactoring.autoSave,omitempty"`
	FilesRestoreUndoStack        bool     `json:"files.restoreUndoStack,omitempty"`
	FilesSaveConflictResolution  string   `json:"files.saveConflictResolution,omitempty"`
	FilesSimpleDialogEnable      bool     `json:"files.simpleDialog.enable,omitempty"`
	FilesTrimFinalNewlines       bool     `json:"files.trimFinalNewlines,omitempty"`
	FilesTrimTrailingWhitespace  bool     `json:"files.trimTrailingWhitespace,omitempty"`
	FilesWatcherExclude          struct {
		GitObjects      bool `json:"**/.git/objects/**,omitempty"`
		GitSubtreeCache bool `json:"**/.git/subtree-cache/**,omitempty"`
		NodeModules     bool `json:"**/node_modules/*/**,omitempty"`
		HgStore         bool `json:"**/.hg/store/**,omitempty"`
	} `json:"files.watcherExclude,omitempty"`
	FilesWatcherInclude           []any `json:"files.watcherInclude,omitempty"`
	ScreencastModeFontSize        int   `json:"screencastMode.fontSize,omitempty"`
	ScreencastModeKeyboardOptions struct {
		ShowKeys                    bool `json:"showKeys,omitempty"`
		ShowCommands                bool `json:"showCommands,omitempty"`
		ShowCommandGroups           bool `json:"showCommandGroups,omitempty"`
		ShowSingleEditorCursorMoves bool `json:"showSingleEditorCursorMoves,omitempty"`
	} `json:"screencastMode.keyboardOptions,omitempty"`
	ScreencastModeKeyboardOverlayTimeouT int    `json:"screencastMode.keyboardOverlayTimeout,omitempty"`
	ScreencastModeMouseIndicatorColor    string `json:"screencastMode.mouseIndicatorColor,omitempty"`
	ScreencastModeMouseIndicatorSize     int    `json:"screencastMode.mouseIndicatorSize,omitempty"`
	ScreencastModeVerticalOffset         int    `json:"screencastMode.verticalOffset,omitempty"`
	ZenModeCenterLayout                  bool   `json:"zenMode.centerLayout,omitempty"`
	ZenModeFullScreen                    bool   `json:"zenMode.fullScreen,omitempty"`
	ZenModeHideActivityBar               bool   `json:"zenMode.hideActivityBar,omitempty"`
	ZenModeHideLineNumbers               bool   `json:"zenMode.hideLineNumbers,omitempty"`
	ZenModeHideStatusBar                 bool   `json:"zenMode.hideStatusBar,omitempty"`
	ZenModeHideTabs                      bool   `json:"zenMode.hideTabs,omitempty"`
	ZenModeRestore                       bool   `json:"zenMode.restore,omitempty"`
	ZenModeSilentNotifications           bool   `json:"zenMode.silentNotifications,omitempty"`
	ExplorerAutoReveal                   bool   `json:"explorer.autoReveal,omitempty"`
	ExplorerAutoRevealExclude            struct {
		NodeModules     bool `json:"**/node_modules,omitempty"`
		BowerComponents bool `json:"**/bower_components,omitempty"`
	} `json:"explorer.autoRevealExclude,omitempty"`
	ExplorerCompactFolders               bool   `json:"explorer.compactFolders,omitempty"`
	ExplorerConfirmDelete                bool   `json:"explorer.confirmDelete,omitempty"`
	ExplorerConfirmDragAndDrop           bool   `json:"explorer.confirmDragAndDrop,omitempty"`
	ExplorerConfirmUndo                  string `json:"explorer.confirmUndo,omitempty"`
	ExplorerCopyRelativePathSeparator    string `json:"explorer.copyRelativePathSeparator,omitempty"`
	ExplorerDecorationsBadges            bool   `json:"explorer.decorations.badges,omitempty"`
	ExplorerDecorationsColors            bool   `json:"explorer.decorations.colors,omitempty"`
	ExplorerEnableDragAndDrop            bool   `json:"explorer.enableDragAndDrop,omitempty"`
	ExplorerEnableUndo                   bool   `json:"explorer.enableUndo,omitempty"`
	ExplorerExcludeGitIgnore             bool   `json:"explorer.excludeGitIgnore,omitempty"`
	ExplorerExpandSingleFolderWorkspaces bool   `json:"explorer.expandSingleFolderWorkspaces,omitempty"`
	ExplorerFileNestingEnabled           bool   `json:"explorer.fileNesting.enabled,omitempty"`
	ExplorerFileNestingExpand            bool   `json:"explorer.fileNesting.expand,omitempty"`
	ExplorerFileNestingPatterns          struct {
		Ts           string `json:"*.ts,omitempty"`
		Js           string `json:"*.js,omitempty"`
		Jsx          string `json:"*.jsx,omitempty"`
		Tsx          string `json:"*.tsx,omitempty"`
		TsconfigJSON string `json:"tsconfig.json,omitempty"`
		PackageJSON  string `json:"package.json,omitempty"`
	} `json:"explorer.fileNesting.patterns,omitempty"`
	ExplorerIncrementalNaming             string `json:"explorer.incrementalNaming,omitempty"`
	ExplorerOpenEditorsMinVisible         int    `json:"explorer.openEditors.minVisible,omitempty"`
	ExplorerOpenEditorsSortOrder          string `json:"explorer.openEditors.sortOrder,omitempty"`
	ExplorerOpenEditorsVisible            int    `json:"explorer.openEditors.visible,omitempty"`
	ExplorerSortOrder                     string `json:"explorer.sortOrder,omitempty"`
	ExplorerSortOrderLexicographicOptions string `json:"explorer.sortOrderLexicographicOptions,omitempty"`
	SearchActionsPosition                 string `json:"search.actionsPosition,omitempty"`
	SearchCollapseResults                 string `json:"search.collapseResults,omitempty"`
	SearchDecorationsBadges               bool   `json:"search.decorations.badges,omitempty"`
	SearchDecorationsColors               bool   `json:"search.decorations.colors,omitempty"`
	SearchDefaultViewMode                 string `json:"search.defaultViewMode,omitempty"`
	SearchExclude                         struct {
		NodeModules     bool `json:"**/node_modules,omitempty"`
		BowerComponents bool `json:"**/bower_components,omitempty"`
		CodeSearch      bool `json:"**/*.code-search,omitempty"`
	} `json:"search.exclude,omitempty"`
	SearchFollowSymlinks                            bool   `json:"search.followSymlinks,omitempty"`
	SearchGlobalFindClipboard                       bool   `json:"search.globalFindClipboard,omitempty"`
	SearchMode                                      string `json:"search.mode,omitempty"`
	SearchQuickOpenHistoryFilterSortOrder           string `json:"search.quickOpen.history.filterSortOrder,omitempty"`
	SearchQuickOpenIncludeHistory                   bool   `json:"search.quickOpen.includeHistory,omitempty"`
	SearchQuickOpenIncludeSymbols                   bool   `json:"search.quickOpen.includeSymbols,omitempty"`
	SearchSearchEditorDefaultNumberOfContextLines   int    `json:"search.searchEditor.defaultNumberOfContextLines,omitempty"`
	SearchSearchEditorDoubleClickBehaviour          string `json:"search.searchEditor.doubleClickBehaviour,omitempty"`
	SearchSearchEditorReusePriorSearchConfiguration bool   `json:"search.searchEditor.reusePriorSearchConfiguration,omitempty"`
	SearchSearchOnType                              bool   `json:"search.searchOnType,omitempty"`
	SearchSearchOnTypeDebouncePeriod                int    `json:"search.searchOnTypeDebouncePeriod,omitempty"`
	SearchSeedOnFocus                               bool   `json:"search.seedOnFocus,omitempty"`
	SearchSeedWithNearestWord                       bool   `json:"search.seedWithNearestWord,omitempty"`
	SearchShowLineNumbers                           bool   `json:"search.showLineNumbers,omitempty"`
	SearchSmartCase                                 bool   `json:"search.smartCase,omitempty"`
	SearchSortOrder                                 string `json:"search.sortOrder,omitempty"`
	SearchUseGlobalIgnoreFiles                      bool   `json:"search.useGlobalIgnoreFiles,omitempty"`
	SearchUseIgnoreFiles                            bool   `json:"search.useIgnoreFiles,omitempty"`
	SearchUseParentIgnoreFiles                      bool   `json:"search.useParentIgnoreFiles,omitempty"`
	SearchUseReplacePreview                         bool   `json:"search.useReplacePreview,omitempty"`
	HTTPProxy                                       string `json:"http.proxy,omitempty"`
	HTTPProxyAuthorization                          any    `json:"http.proxyAuthorization,omitempty"`
	HTTPProxyKerberosServicePrincipal               string `json:"http.proxyKerberosServicePrincipal,omitempty"`
	HTTPProxyStrictSSL                              bool   `json:"http.proxyStrictSSL,omitempty"`
	HTTPProxySupport                                string `json:"http.proxySupport,omitempty"`
	HTTPSystemCertificates                          bool   `json:"http.systemCertificates,omitempty"`
	KeyboardDispatch                                string `json:"keyboard.dispatch,omitempty"`
	KeyboardMapAltGrToCtrlAlt                       bool   `json:"keyboard.mapAltGrToCtrlAlt,omitempty"`
	KeyboardTouchbarEnabled                         bool   `json:"keyboard.touchbar.enabled,omitempty"`
	KeyboardTouchbarIgnored                         []any  `json:"keyboard.touchbar.ignored,omitempty"`
	UpdateEnableWindowsBackgroundUpdates            bool   `json:"update.enableWindowsBackgroundUpdates,omitempty"`
	UpdateMode                                      string `json:"update.mode,omitempty"`
	UpdateShowReleaseNotes                          bool   `json:"update.showReleaseNotes,omitempty"`
	DebugAllowBreakpointsEverywhere                 bool   `json:"debug.allowBreakpointsEverywhere,omitempty"`
	DebugAutoExpandLazyVariables                    bool   `json:"debug.autoExpandLazyVariables,omitempty"`
	DebugConfirmOnExit                              string `json:"debug.confirmOnExit,omitempty"`
	DebugConsoleAcceptSuggestionOnEnter             string `json:"debug.console.acceptSuggestionOnEnter,omitempty"`
	DebugConsoleCloseOnEnd                          bool   `json:"debug.console.closeOnEnd,omitempty"`
	DebugConsoleCollapseIdenticalLines              bool   `json:"debug.console.collapseIdenticalLines,omitempty"`
	DebugConsoleFontFamily                          string `json:"debug.console.fontFamily,omitempty"`
	DebugConsoleFontSize                            int    `json:"debug.console.fontSize,omitempty"`
	DebugConsoleHistorySuggestions                  bool   `json:"debug.console.historySuggestions,omitempty"`
	DebugConsoleLineHeight                          int    `json:"debug.console.lineHeight,omitempty"`
	DebugConsoleWordWrap                            bool   `json:"debug.console.wordWrap,omitempty"`
	DebugDisassemblyViewShowSourceCode              bool   `json:"debug.disassemblyView.showSourceCode,omitempty"`
	DebugEnableStatusBarColor                       bool   `json:"debug.enableStatusBarColor,omitempty"`
	DebugFocusEditorOnBreak                         bool   `json:"debug.focusEditorOnBreak,omitempty"`
	DebugFocusWindowOnBreak                         bool   `json:"debug.focusWindowOnBreak,omitempty"`
	DebugInlineValues                               string `json:"debug.inlineValues,omitempty"`
	DebugInternalConsoleOptions                     string `json:"debug.internalConsoleOptions,omitempty"`
	DebugOnTaskErrors                               string `json:"debug.onTaskErrors,omitempty"`
	DebugOpenDebug                                  string `json:"debug.openDebug,omitempty"`
	DebugOpenExplorerOnEnd                          bool   `json:"debug.openExplorerOnEnd,omitempty"`
	DebugSaveBeforeStart                            string `json:"debug.saveBeforeStart,omitempty"`
	DebugShowBreakpointsInOverviewRuler             bool   `json:"debug.showBreakpointsInOverviewRuler,omitempty"`
	DebugShowInlineBreakpointCandidates             bool   `json:"debug.showInlineBreakpointCandidates,omitempty"`
	DebugShowInStatusBar                            string `json:"debug.showInStatusBar,omitempty"`
	DebugShowSubSessionsInToolBar                   bool   `json:"debug.showSubSessionsInToolBar,omitempty"`
	DebugTerminalClearBeforeReusing                 bool   `json:"debug.terminal.clearBeforeReusing,omitempty"`
	DebugToolBarLocation                            string `json:"debug.toolBarLocation,omitempty"`
	Launch                                          struct {
		Configurations []any `json:"configurations,omitempty"`
		Compounds      []any `json:"compounds,omitempty"`
	} `json:"launch,omitempty"`
	HTMLAutoClosingTags                                                         bool     `json:"html.autoClosingTags,omitempty"`
	HTMLAutoCreateQuotes                                                        bool     `json:"html.autoCreateQuotes,omitempty"`
	HTMLCompletionAttributeDefaultValue                                         string   `json:"html.completion.attributeDefaultValue,omitempty"`
	HTMLCustomData                                                              []any    `json:"html.customData,omitempty"`
	HTMLFormatContentUnformatted                                                string   `json:"html.format.contentUnformatted,omitempty"`
	HTMLFormatEnable                                                            bool     `json:"html.format.enable,omitempty"`
	HTMLFormatExtraLiners                                                       string   `json:"html.format.extraLiners,omitempty"`
	HTMLFormatIndentHandlebars                                                  bool     `json:"html.format.indentHandlebars,omitempty"`
	HTMLFormatIndentInnerHTML                                                   bool     `json:"html.format.indentInnerHtml,omitempty"`
	HTMLFormatMaxPreserveNewLines                                               any      `json:"html.format.maxPreserveNewLines,omitempty"`
	HTMLFormatPreserveNewLines                                                  bool     `json:"html.format.preserveNewLines,omitempty"`
	HTMLFormatTemplating                                                        bool     `json:"html.format.templating,omitempty"`
	HTMLFormatUnformatted                                                       string   `json:"html.format.unformatted,omitempty"`
	HTMLFormatUnformattedContentDelimiter                                       string   `json:"html.format.unformattedContentDelimiter,omitempty"`
	HTMLFormatWrapAttributes                                                    string   `json:"html.format.wrapAttributes,omitempty"`
	HTMLFormatWrapAttributesIndentSize                                          any      `json:"html.format.wrapAttributesIndentSize,omitempty"`
	HTMLFormatWrapLineLength                                                    int      `json:"html.format.wrapLineLength,omitempty"`
	HTMLHoverDocumentation                                                      bool     `json:"html.hover.documentation,omitempty"`
	HTMLHoverReferences                                                         bool     `json:"html.hover.references,omitempty"`
	HTMLSuggestHTML5                                                            bool     `json:"html.suggest.html5,omitempty"`
	HTMLTraceServer                                                             string   `json:"html.trace.server,omitempty"`
	HTMLValidateScripts                                                         bool     `json:"html.validate.scripts,omitempty"`
	HTMLValidateStyles                                                          bool     `json:"html.validate.styles,omitempty"`
	JSONFormatEnable                                                            bool     `json:"json.format.enable,omitempty"`
	JSONFormatKeepLines                                                         bool     `json:"json.format.keepLines,omitempty"`
	JSONMaxItemsComputed                                                        int      `json:"json.maxItemsComputed,omitempty"`
	JSONSchemaDownloadEnable                                                    bool     `json:"json.schemaDownload.enable,omitempty"`
	JSONSchemas                                                                 []any    `json:"json.schemas,omitempty"`
	JSONTraceServer                                                             string   `json:"json.trace.server,omitempty"`
	JSONValidateEnable                                                          bool     `json:"json.validate.enable,omitempty"`
	MarkdownCopyFilesDestination                                                struct{} `json:"markdown.copyFiles.destination,omitempty"`
	MarkdownCopyFilesOverwriteBehavior                                          string   `json:"markdown.copyFiles.overwriteBehavior,omitempty"`
	MarkdownEditorDropCopyIntoWorkspace                                         string   `json:"markdown.editor.drop.copyIntoWorkspace,omitempty"`
	MarkdownEditorDropEnabled                                                   bool     `json:"markdown.editor.drop.enabled,omitempty"`
	MarkdownEditorFilePasteCopyIntoWorkspace                                    string   `json:"markdown.editor.filePaste.copyIntoWorkspace,omitempty"`
	MarkdownEditorFilePasteEnabled                                              bool     `json:"markdown.editor.filePaste.enabled,omitempty"`
	MarkdownEditorPasteURLAsFormattedLinkEnabled                                string   `json:"markdown.editor.pasteUrlAsFormattedLink.enabled,omitempty"`
	MarkdownLinksOpenLocation                                                   string   `json:"markdown.links.openLocation,omitempty"`
	MarkdownOccurrencesHighlightEnabled                                         bool     `json:"markdown.occurrencesHighlight.enabled,omitempty"`
	MarkdownPreferredMdPathExtensionStyle                                       string   `json:"markdown.preferredMdPathExtensionStyle,omitempty"`
	MarkdownPreviewBreaks                                                       bool     `json:"markdown.preview.breaks,omitempty"`
	MarkdownPreviewDoubleClickToSwitchToEditor                                  bool     `json:"markdown.preview.doubleClickToSwitchToEditor,omitempty"`
	MarkdownPreviewFontFamily                                                   string   `json:"markdown.preview.fontFamily,omitempty"`
	MarkdownPreviewFontSize                                                     int      `json:"markdown.preview.fontSize,omitempty"`
	MarkdownPreviewLineHeight                                                   float64  `json:"markdown.preview.lineHeight,omitempty"`
	MarkdownPreviewLinkify                                                      bool     `json:"markdown.preview.linkify,omitempty"`
	MarkdownPreviewMarkEditorSelection                                          bool     `json:"markdown.preview.markEditorSelection,omitempty"`
	MarkdownPreviewOpenMarkdownLinks                                            string   `json:"markdown.preview.openMarkdownLinks,omitempty"`
	MarkdownPreviewScrollEditorWithPreview                                      bool     `json:"markdown.preview.scrollEditorWithPreview,omitempty"`
	MarkdownPreviewScrollPreviewWithEditor                                      bool     `json:"markdown.preview.scrollPreviewWithEditor,omitempty"`
	MarkdownPreviewTypographer                                                  bool     `json:"markdown.preview.typographer,omitempty"`
	MarkdownStyles                                                              []any    `json:"markdown.styles,omitempty"`
	MarkdownSuggestPathsEnabled                                                 bool     `json:"markdown.suggest.paths.enabled,omitempty"`
	MarkdownSuggestPathsIncludeWorkspaceHeaderCompletions                       string   `json:"markdown.suggest.paths.includeWorkspaceHeaderCompletions,omitempty"`
	MarkdownTraceExtension                                                      string   `json:"markdown.trace.extension,omitempty"`
	MarkdownUpdateLinksOnFileMoveEnabled                                        string   `json:"markdown.updateLinksOnFileMove.enabled,omitempty"`
	MarkdownUpdateLinksOnFileMoveEnableForDirectories                           bool     `json:"markdown.updateLinksOnFileMove.enableForDirectories,omitempty"`
	MarkdownUpdateLinksOnFileMoveInclude                                        []string `json:"markdown.updateLinksOnFileMove.include,omitempty"`
	MarkdownValidateDuplicateLinkDefinitionsEnabled                             string   `json:"markdown.validate.duplicateLinkDefinitions.enabled,omitempty"`
	MarkdownValidateEnabled                                                     bool     `json:"markdown.validate.enabled,omitempty"`
	MarkdownValidateFileLinksEnabled                                            string   `json:"markdown.validate.fileLinks.enabled,omitempty"`
	MarkdownValidateFileLinksMarkdownFragmentLinks                              string   `json:"markdown.validate.fileLinks.markdownFragmentLinks,omitempty"`
	MarkdownValidateFragmentLinksEnabled                                        string   `json:"markdown.validate.fragmentLinks.enabled,omitempty"`
	MarkdownValidateIgnoredLinks                                                []any    `json:"markdown.validate.ignoredLinks,omitempty"`
	MarkdownValidateReferenceLinksEnabled                                       string   `json:"markdown.validate.referenceLinks.enabled,omitempty"`
	MarkdownValidateUnusedLinkDefinitionsEnabled                                string   `json:"markdown.validate.unusedLinkDefinitions.enabled,omitempty"`
	PhpSuggestBasic                                                             bool     `json:"php.suggest.basic,omitempty"`
	PhpValidateEnable                                                           bool     `json:"php.validate.enable,omitempty"`
	PhpValidateExecutablePath                                                   any      `json:"php.validate.executablePath,omitempty"`
	PhpValidateRun                                                              string   `json:"php.validate.run,omitempty"`
	JavascriptAutoClosingTags                                                   bool     `json:"javascript.autoClosingTags,omitempty"`
	JavascriptFormatEnable                                                      bool     `json:"javascript.format.enable,omitempty"`
	JavascriptFormatInsertSpaceAfterCommaDelimiter                              bool     `json:"javascript.format.insertSpaceAfterCommaDelimiter,omitempty"`
	JavascriptFormatInsertSpaceAfterConstructor                                 bool     `json:"javascript.format.insertSpaceAfterConstructor,omitempty"`
	JavascriptFormatInsertSpaceAfterFunctionKeywordForAnonymousFunctions        bool     `json:"javascript.format.insertSpaceAfterFunctionKeywordForAnonymousFunctions,omitempty"`
	JavascriptFormatInsertSpaceAfterKeywordsInControlFlowStatements             bool     `json:"javascript.format.insertSpaceAfterKeywordsInControlFlowStatements,omitempty"`
	JavascriptFormatInsertSpaceAfterOpeningAndBeforeClosingEmptyBraces          bool     `json:"javascript.format.insertSpaceAfterOpeningAndBeforeClosingEmptyBraces,omitempty"`
	JavascriptFormatInsertSpaceAfterOpeningAndBeforeClosingJsxExpressionBraces  bool     `json:"javascript.format.insertSpaceAfterOpeningAndBeforeClosingJsxExpressionBraces,omitempty"`
	JavascriptFormatInsertSpaceAfterOpeningAndBeforeClosingNonemptyBraces       bool     `json:"javascript.format.insertSpaceAfterOpeningAndBeforeClosingNonemptyBraces,omitempty"`
	JavascriptFormatInsertSpaceAfterOpeningAndBeforeClosingNonemptyBrackets     bool     `json:"javascript.format.insertSpaceAfterOpeningAndBeforeClosingNonemptyBrackets,omitempty"`
	JavascriptFormatInsertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis  bool     `json:"javascript.format.insertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis,omitempty"`
	JavascriptFormatInsertSpaceAfterOpeningAndBeforeClosingTemplateStringBraces bool     `json:"javascript.format.insertSpaceAfterOpeningAndBeforeClosingTemplateStringBraces,omitempty"`
	JavascriptFormatInsertSpaceAfterSemicolonInForStatements                    bool     `json:"javascript.format.insertSpaceAfterSemicolonInForStatements,omitempty"`
	JavascriptFormatInsertSpaceBeforeAndAfterBinaryOperators                    bool     `json:"javascript.format.insertSpaceBeforeAndAfterBinaryOperators,omitempty"`
	JavascriptFormatInsertSpaceBeforeFunctionParenthesis                        bool     `json:"javascript.format.insertSpaceBeforeFunctionParenthesis,omitempty"`
	JavascriptFormatPlaceOpenBraceOnNewLineForControlBlocks                     bool     `json:"javascript.format.placeOpenBraceOnNewLineForControlBlocks,omitempty"`
	JavascriptFormatPlaceOpenBraceOnNewLineForFunctions                         bool     `json:"javascript.format.placeOpenBraceOnNewLineForFunctions,omitempty"`
	JavascriptFormatSemicolons                                                  string   `json:"javascript.format.semicolons,omitempty"`
	JavascriptInlayHintsEnumMemberValuesEnabled                                 bool     `json:"javascript.inlayHints.enumMemberValues.enabled,omitempty"`
	JavascriptInlayHintsFunctionLikeReturnTypesEnabled                          bool     `json:"javascript.inlayHints.functionLikeReturnTypes.enabled,omitempty"`
	JavascriptInlayHintsParameterNamesEnabled                                   string   `json:"javascript.inlayHints.parameterNames.enabled,omitempty"`
	JavascriptInlayHintsParameterNamesSuppressWhenArgumentMatchesName           bool     `json:"javascript.inlayHints.parameterNames.suppressWhenArgumentMatchesName,omitempty"`
	JavascriptInlayHintsParameterTypesEnabled                                   bool     `json:"javascript.inlayHints.parameterTypes.enabled,omitempty"`
	JavascriptInlayHintsPropertyDeclarationTypesEnabled                         bool     `json:"javascript.inlayHints.propertyDeclarationTypes.enabled,omitempty"`
	JavascriptInlayHintsVariableTypesEnabled                                    bool     `json:"javascript.inlayHints.variableTypes.enabled,omitempty"`
	JavascriptInlayHintsVariableTypesSuppressWhenTypeMatchesName                bool     `json:"javascript.inlayHints.variableTypes.suppressWhenTypeMatchesName,omitempty"`
	JavascriptPreferencesAutoImportFileExcludePatterns                          []any    `json:"javascript.preferences.autoImportFileExcludePatterns,omitempty"`
	JavascriptPreferencesImportModuleSpecifier                                  string   `json:"javascript.preferences.importModuleSpecifier,omitempty"`
	JavascriptPreferencesImportModuleSpecifierEnding                            string   `json:"javascript.preferences.importModuleSpecifierEnding,omitempty"`
	JavascriptPreferencesJsxAttributeCompletionStyle                            string   `json:"javascript.preferences.jsxAttributeCompletionStyle,omitempty"`
	JavascriptPreferencesQuoteStyle                                             string   `json:"javascript.preferences.quoteStyle,omitempty"`
	JavascriptPreferencesRenameMatchingJsxTags                                  bool     `json:"javascript.preferences.renameMatchingJsxTags,omitempty"`
	JavascriptPreferencesRenameShorthandProperties                              bool     `json:"javascript.preferences.renameShorthandProperties,omitempty"`
	JavascriptPreferencesUseAliasesForRenames                                   bool     `json:"javascript.preferences.useAliasesForRenames,omitempty"`
	JavascriptPreferGoToSourceDefinition                                        bool     `json:"javascript.preferGoToSourceDefinition,omitempty"`
	JavascriptReferencesCodeLensEnabled                                         bool     `json:"javascript.referencesCodeLens.enabled,omitempty"`
	JavascriptReferencesCodeLensShowOnAllFunctions                              bool     `json:"javascript.referencesCodeLens.showOnAllFunctions,omitempty"`
	JavascriptSuggestAutoImports                                                bool     `json:"javascript.suggest.autoImports,omitempty"`
	JavascriptSuggestClassMemberSnippetsEnabled                                 bool     `json:"javascript.suggest.classMemberSnippets.enabled,omitempty"`
	JavascriptSuggestCompleteFunctionCalls                                      bool     `json:"javascript.suggest.completeFunctionCalls,omitempty"`
	JavascriptSuggestCompleteJSDocs                                             bool     `json:"javascript.suggest.completeJSDocs,omitempty"`
	JavascriptSuggestEnabled                                                    bool     `json:"javascript.suggest.enabled,omitempty"`
	JavascriptSuggestIncludeAutomaticOptionalChainCompletions                   bool     `json:"javascript.suggest.includeAutomaticOptionalChainCompletions,omitempty"`
	JavascriptSuggestIncludeCompletionsForImportStatements                      bool     `json:"javascript.suggest.includeCompletionsForImportStatements,omitempty"`
	JavascriptSuggestJsdocGenerateReturns                                       bool     `json:"javascript.suggest.jsdoc.generateReturns,omitempty"`
	JavascriptSuggestNames                                                      bool     `json:"javascript.suggest.names,omitempty"`
	JavascriptSuggestPaths                                                      bool     `json:"javascript.suggest.paths,omitempty"`
	JavascriptSuggestionActionsEnabled                                          bool     `json:"javascript.suggestionActions.enabled,omitempty"`
	JavascriptUpdateImportsOnFileMoveEnabled                                    string   `json:"javascript.updateImportsOnFileMove.enabled,omitempty"`
	JavascriptValidateEnable                                                    bool     `json:"javascript.validate.enable,omitempty"`
	JsTsImplicitProjectConfigCheckJs                                            bool     `json:"js/ts.implicitProjectConfig.checkJs,omitempty"`
	JsTsImplicitProjectConfigModule                                             string   `json:"js/ts.implicitProjectConfig.module,omitempty"`
	JsTsImplicitProjectConfigStrictFunctionTypes                                bool     `json:"js/ts.implicitProjectConfig.strictFunctionTypes,omitempty"`
	JsTsImplicitProjectConfigStrictNullChecks                                   bool     `json:"js/ts.implicitProjectConfig.strictNullChecks,omitempty"`
	JsTsImplicitProjectConfigTarget                                             string   `json:"js/ts.implicitProjectConfig.target,omitempty"`
	TypescriptAutoClosingTags                                                   bool     `json:"typescript.autoClosingTags,omitempty"`
	TypescriptCheckNpmIsInstalled                                               bool     `json:"typescript.check.npmIsInstalled,omitempty"`
	TypescriptDisableAutomaticTypeAcquisition                                   bool     `json:"typescript.disableAutomaticTypeAcquisition,omitempty"`
	TypescriptEnablePromptUseWorkspaceTsdk                                      bool     `json:"typescript.enablePromptUseWorkspaceTsdk,omitempty"`
	TypescriptFormatEnable                                                      bool     `json:"typescript.format.enable,omitempty"`
	TypescriptFormatIndentSwitchCase                                            bool     `json:"typescript.format.indentSwitchCase,omitempty"`
	TypescriptFormatInsertSpaceAfterCommaDelimiter                              bool     `json:"typescript.format.insertSpaceAfterCommaDelimiter,omitempty"`
	TypescriptFormatInsertSpaceAfterConstructor                                 bool     `json:"typescript.format.insertSpaceAfterConstructor,omitempty"`
	TypescriptFormatInsertSpaceAfterFunctionKeywordForAnonymousFunctions        bool     `json:"typescript.format.insertSpaceAfterFunctionKeywordForAnonymousFunctions,omitempty"`
	TypescriptFormatInsertSpaceAfterKeywordsInControlFlowStatements             bool     `json:"typescript.format.insertSpaceAfterKeywordsInControlFlowStatements,omitempty"`
	TypescriptFormatInsertSpaceAfterOpeningAndBeforeClosingEmptyBraces          bool     `json:"typescript.format.insertSpaceAfterOpeningAndBeforeClosingEmptyBraces,omitempty"`
	TypescriptFormatInsertSpaceAfterOpeningAndBeforeClosingJsxExpressionBraces  bool     `json:"typescript.format.insertSpaceAfterOpeningAndBeforeClosingJsxExpressionBraces,omitempty"`
	TypescriptFormatInsertSpaceAfterOpeningAndBeforeClosingNonemptyBraces       bool     `json:"typescript.format.insertSpaceAfterOpeningAndBeforeClosingNonemptyBraces,omitempty"`
	TypescriptFormatInsertSpaceAfterOpeningAndBeforeClosingNonemptyBrackets     bool     `json:"typescript.format.insertSpaceAfterOpeningAndBeforeClosingNonemptyBrackets,omitempty"`
	TypescriptFormatInsertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis  bool     `json:"typescript.format.insertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis,omitempty"`
	TypescriptFormatInsertSpaceAfterOpeningAndBeforeClosingTemplateStringBraces bool     `json:"typescript.format.insertSpaceAfterOpeningAndBeforeClosingTemplateStringBraces,omitempty"`
	TypescriptFormatInsertSpaceAfterSemicolonInForStatements                    bool     `json:"typescript.format.insertSpaceAfterSemicolonInForStatements,omitempty"`
	TypescriptFormatInsertSpaceAfterTypeAssertion                               bool     `json:"typescript.format.insertSpaceAfterTypeAssertion,omitempty"`
	TypescriptFormatInsertSpaceBeforeAndAfterBinaryOperators                    bool     `json:"typescript.format.insertSpaceBeforeAndAfterBinaryOperators,omitempty"`
	TypescriptFormatInsertSpaceBeforeFunctionParenthesis                        bool     `json:"typescript.format.insertSpaceBeforeFunctionParenthesis,omitempty"`
	TypescriptFormatPlaceOpenBraceOnNewLineForControlBlocks                     bool     `json:"typescript.format.placeOpenBraceOnNewLineForControlBlocks,omitempty"`
	TypescriptFormatPlaceOpenBraceOnNewLineForFunctions                         bool     `json:"typescript.format.placeOpenBraceOnNewLineForFunctions,omitempty"`
	TypescriptFormatSemicolons                                                  string   `json:"typescript.format.semicolons,omitempty"`
	TypescriptImplementationsCodeLensEnabled                                    bool     `json:"typescript.implementationsCodeLens.enabled,omitempty"`
	TypescriptInlayHintsEnumMemberValuesEnabled                                 bool     `json:"typescript.inlayHints.enumMemberValues.enabled,omitempty"`
	TypescriptInlayHintsFunctionLikeReturnTypesEnabled                          bool     `json:"typescript.inlayHints.functionLikeReturnTypes.enabled,omitempty"`
	TypescriptInlayHintsParameterNamesEnabled                                   string   `json:"typescript.inlayHints.parameterNames.enabled,omitempty"`
	TypescriptInlayHintsParameterNamesSuppressWhenArgumentMatchesName           bool     `json:"typescript.inlayHints.parameterNames.suppressWhenArgumentMatchesName,omitempty"`
	TypescriptInlayHintsParameterTypesEnabled                                   bool     `json:"typescript.inlayHints.parameterTypes.enabled,omitempty"`
	TypescriptInlayHintsPropertyDeclarationTypesEnabled                         bool     `json:"typescript.inlayHints.propertyDeclarationTypes.enabled,omitempty"`
	TypescriptInlayHintsVariableTypesEnabled                                    bool     `json:"typescript.inlayHints.variableTypes.enabled,omitempty"`
	TypescriptInlayHintsVariableTypesSuppressWhenTypeMatchesName                bool     `json:"typescript.inlayHints.variableTypes.suppressWhenTypeMatchesName,omitempty"`
	TypescriptLocale                                                            string   `json:"typescript.locale,omitempty"`
	TypescriptNpm                                                               string   `json:"typescript.npm,omitempty"`
	TypescriptPreferencesAutoImportFileExcludePatterns                          []any    `json:"typescript.preferences.autoImportFileExcludePatterns,omitempty"`
	TypescriptPreferencesImportModuleSpecifier                                  string   `json:"typescript.preferences.importModuleSpecifier,omitempty"`
	TypescriptPreferencesImportModuleSpecifierEnding                            string   `json:"typescript.preferences.importModuleSpecifierEnding,omitempty"`
	TypescriptPreferencesIncludePackageJSONAutoImports                          string   `json:"typescript.preferences.includePackageJsonAutoImports,omitempty"`
	TypescriptPreferencesJsxAttributeCompletionStyle                            string   `json:"typescript.preferences.jsxAttributeCompletionStyle,omitempty"`
	TypescriptPreferencesQuoteStyle                                             string   `json:"typescript.preferences.quoteStyle,omitempty"`
	TypescriptPreferencesRenameMatchingJsxTags                                  bool     `json:"typescript.preferences.renameMatchingJsxTags,omitempty"`
	TypescriptPreferencesUseAliasesForRenames                                   bool     `json:"typescript.preferences.useAliasesForRenames,omitempty"`
	TypescriptPreferGoToSourceDefinition                                        bool     `json:"typescript.preferGoToSourceDefinition,omitempty"`
	TypescriptReferencesCodeLensEnabled                                         bool     `json:"typescript.referencesCodeLens.enabled,omitempty"`
	TypescriptReferencesCodeLensShowOnAllFunctions                              bool     `json:"typescript.referencesCodeLens.showOnAllFunctions,omitempty"`
	TypescriptReportStyleChecksAsWarnings                                       bool     `json:"typescript.reportStyleChecksAsWarnings,omitempty"`
	TypescriptSuggestAutoImports                                                bool     `json:"typescript.suggest.autoImports,omitempty"`
	TypescriptSuggestClassMemberSnippetsEnabled                                 bool     `json:"typescript.suggest.classMemberSnippets.enabled,omitempty"`
	TypescriptSuggestCompleteFunctionCalls                                      bool     `json:"typescript.suggest.completeFunctionCalls,omitempty"`
	TypescriptSuggestCompleteJSDocs                                             bool     `json:"typescript.suggest.completeJSDocs,omitempty"`
	TypescriptSuggestEnabled                                                    bool     `json:"typescript.suggest.enabled,omitempty"`
	TypescriptSuggestIncludeAutomaticOptionalChainCompletions                   bool     `json:"typescript.suggest.includeAutomaticOptionalChainCompletions,omitempty"`
	TypescriptSuggestIncludeCompletionsForImportStatements                      bool     `json:"typescript.suggest.includeCompletionsForImportStatements,omitempty"`
	TypescriptSuggestJsdocGenerateReturns                                       bool     `json:"typescript.suggest.jsdoc.generateReturns,omitempty"`
	TypescriptSuggestObjectLiteralMethodSnippetsEnabled                         bool     `json:"typescript.suggest.objectLiteralMethodSnippets.enabled,omitempty"`
	TypescriptSuggestPaths                                                      bool     `json:"typescript.suggest.paths,omitempty"`
	TypescriptSuggestionActionsEnabled                                          bool     `json:"typescript.suggestionActions.enabled,omitempty"`
	TypescriptTscAutoDetect                                                     string   `json:"typescript.tsc.autoDetect,omitempty"`
	TypescriptTsdk                                                              string   `json:"typescript.tsdk,omitempty"`
	TypescriptTsserverEnableTracing                                             bool     `json:"typescript.tsserver.enableTracing,omitempty"`
	TypescriptTsserverLog                                                       string   `json:"typescript.tsserver.log,omitempty"`
	TypescriptTsserverMaxTsServerMemory                                         int      `json:"typescript.tsserver.maxTsServerMemory,omitempty"`
	TypescriptTsserverPluginPaths                                               []any    `json:"typescript.tsserver.pluginPaths,omitempty"`
	TypescriptTsserverUseSyntaxServer                                           string   `json:"typescript.tsserver.useSyntaxServer,omitempty"`
	TypescriptTsserverWatchOptions                                              struct{} `json:"typescript.tsserver.watchOptions,omitempty"`
	TypescriptTsserverWebProjectWideIntellisenseEnabled                         bool     `json:"typescript.tsserver.web.projectWideIntellisense.enabled,omitempty"`
	TypescriptTsserverWebProjectWideIntellisenseSuppressSemanticErrors          bool     `json:"typescript.tsserver.web.projectWideIntellisense.suppressSemanticErrors,omitempty"`
	TypescriptUpdateImportsOnFileMoveEnabled                                    string   `json:"typescript.updateImportsOnFileMove.enabled,omitempty"`
	TypescriptValidateEnable                                                    bool     `json:"typescript.validate.enable,omitempty"`
	TypescriptWorkspaceSymbolsScope                                             string   `json:"typescript.workspaceSymbols.scope,omitempty"`
	TestingAlwaysRevealTestOnStateChange                                        bool     `json:"testing.alwaysRevealTestOnStateChange,omitempty"`
	TestingAutomaticallyOpenPeekView                                            string   `json:"testing.automaticallyOpenPeekView,omitempty"`
	TestingAutomaticallyOpenPeekViewDuringAutoRun                               bool     `json:"testing.automaticallyOpenPeekViewDuringAutoRun,omitempty"`
	TestingAutoRunDelay                                                         int      `json:"testing.autoRun.delay,omitempty"`
	TestingCountBadge                                                           string   `json:"testing.countBadge,omitempty"`
	TestingDefaultGutterClickAction                                             string   `json:"testing.defaultGutterClickAction,omitempty"`
	TestingFollowRunningTest                                                    bool     `json:"testing.followRunningTest,omitempty"`
	TestingGutterEnabled                                                        bool     `json:"testing.gutterEnabled,omitempty"`
	TestingOpenTesting                                                          string   `json:"testing.openTesting,omitempty"`
	TestingSaveBeforeTest                                                       bool     `json:"testing.saveBeforeTest,omitempty"`
	TestingShowAllMessages                                                      bool     `json:"testing.showAllMessages,omitempty"`
	CSSCompletionCompletePropertyWithSemicolon                                  bool     `json:"css.completion.completePropertyWithSemicolon,omitempty"`
	CSSCompletionTriggerPropertyValueCompletion                                 bool     `json:"css.completion.triggerPropertyValueCompletion,omitempty"`
	CSSCustomData                                                               []any    `json:"css.customData,omitempty"`
	CSSFormatBraceStyle                                                         string   `json:"css.format.braceStyle,omitempty"`
	CSSFormatEnable                                                             bool     `json:"css.format.enable,omitempty"`
	CSSFormatMaxPreserveNewLines                                                any      `json:"css.format.maxPreserveNewLines,omitempty"`
	CSSFormatNewlineBetweenRules                                                bool     `json:"css.format.newlineBetweenRules,omitempty"`
	CSSFormatNewlineBetweenSelectors                                            bool     `json:"css.format.newlineBetweenSelectors,omitempty"`
	CSSFormatPreserveNewLines                                                   bool     `json:"css.format.preserveNewLines,omitempty"`
	CSSFormatSpaceAroundSelectorSeparator                                       bool     `json:"css.format.spaceAroundSelectorSeparator,omitempty"`
	CSSHoverDocumentation                                                       bool     `json:"css.hover.documentation,omitempty"`
	CSSHoverReferences                                                          bool     `json:"css.hover.references,omitempty"`
	CSSLintArgumentsInColorFunction                                             string   `json:"css.lint.argumentsInColorFunction,omitempty"`
	CSSLintBoxModel                                                             string   `json:"css.lint.boxModel,omitempty"`
	CSSLintCompatibleVendorPrefixes                                             string   `json:"css.lint.compatibleVendorPrefixes,omitempty"`
	CSSLintDuplicateProperties                                                  string   `json:"css.lint.duplicateProperties,omitempty"`
	CSSLintEmptyRules                                                           string   `json:"css.lint.emptyRules,omitempty"`
	CSSLintFloat                                                                string   `json:"css.lint.float,omitempty"`
	CSSLintFontFaceProperties                                                   string   `json:"css.lint.fontFaceProperties,omitempty"`
	CSSLintHexColorLength                                                       string   `json:"css.lint.hexColorLength,omitempty"`
	CSSLintIDSelector                                                           string   `json:"css.lint.idSelector,omitempty"`
	CSSLintIeHack                                                               string   `json:"css.lint.ieHack,omitempty"`
	CSSLintImportant                                                            string   `json:"css.lint.important,omitempty"`
	CSSLintImportStatement                                                      string   `json:"css.lint.importStatement,omitempty"`
	CSSLintPropertyIgnoredDueToDisplay                                          string   `json:"css.lint.propertyIgnoredDueToDisplay,omitempty"`
	CSSLintUniversalSelector                                                    string   `json:"css.lint.universalSelector,omitempty"`
	CSSLintUnknownAtRules                                                       string   `json:"css.lint.unknownAtRules,omitempty"`
	CSSLintUnknownProperties                                                    string   `json:"css.lint.unknownProperties,omitempty"`
	CSSLintUnknownVendorSpecificProperties                                      string   `json:"css.lint.unknownVendorSpecificProperties,omitempty"`
	CSSLintValidProperties                                                      []any    `json:"css.lint.validProperties,omitempty"`
	CSSLintVendorPrefix                                                         string   `json:"css.lint.vendorPrefix,omitempty"`
	CSSLintZeroUnits                                                            string   `json:"css.lint.zeroUnits,omitempty"`
	CSSTraceServer                                                              string   `json:"css.trace.server,omitempty"`
	CSSValidate                                                                 bool     `json:"css.validate,omitempty"`
	LessCompletionCompletePropertyWithSemicolon                                 bool     `json:"less.completion.completePropertyWithSemicolon,omitempty"`
	LessCompletionTriggerPropertyValueCompletion                                bool     `json:"less.completion.triggerPropertyValueCompletion,omitempty"`
	LessFormatBraceStyle                                                        string   `json:"less.format.braceStyle,omitempty"`
	LessFormatEnable                                                            bool     `json:"less.format.enable,omitempty"`
	LessFormatMaxPreserveNewLines                                               any      `json:"less.format.maxPreserveNewLines,omitempty"`
	LessFormatNewlineBetweenRules                                               bool     `json:"less.format.newlineBetweenRules,omitempty"`
	LessFormatNewlineBetweenSelectors                                           bool     `json:"less.format.newlineBetweenSelectors,omitempty"`
	LessFormatPreserveNewLines                                                  bool     `json:"less.format.preserveNewLines,omitempty"`
	LessFormatSpaceAroundSelectorSeparator                                      bool     `json:"less.format.spaceAroundSelectorSeparator,omitempty"`
	LessHoverDocumentation                                                      bool     `json:"less.hover.documentation,omitempty"`
	LessHoverReferences                                                         bool     `json:"less.hover.references,omitempty"`
	LessLintArgumentsInColorFunction                                            string   `json:"less.lint.argumentsInColorFunction,omitempty"`
	LessLintBoxModel                                                            string   `json:"less.lint.boxModel,omitempty"`
	LessLintCompatibleVendorPrefixes                                            string   `json:"less.lint.compatibleVendorPrefixes,omitempty"`
	LessLintDuplicateProperties                                                 string   `json:"less.lint.duplicateProperties,omitempty"`
	LessLintEmptyRules                                                          string   `json:"less.lint.emptyRules,omitempty"`
	LessLintFloat                                                               string   `json:"less.lint.float,omitempty"`
	LessLintFontFaceProperties                                                  string   `json:"less.lint.fontFaceProperties,omitempty"`
	LessLintHexColorLength                                                      string   `json:"less.lint.hexColorLength,omitempty"`
	LessLintIDSelector                                                          string   `json:"less.lint.idSelector,omitempty"`
	LessLintIeHack                                                              string   `json:"less.lint.ieHack,omitempty"`
	LessLintImportant                                                           string   `json:"less.lint.important,omitempty"`
	LessLintImportStatement                                                     string   `json:"less.lint.importStatement,omitempty"`
	LessLintPropertyIgnoredDueToDisplay                                         string   `json:"less.lint.propertyIgnoredDueToDisplay,omitempty"`
	LessLintUniversalSelector                                                   string   `json:"less.lint.universalSelector,omitempty"`
	LessLintUnknownAtRules                                                      string   `json:"less.lint.unknownAtRules,omitempty"`
	LessLintUnknownProperties                                                   string   `json:"less.lint.unknownProperties,omitempty"`
	LessLintUnknownVendorSpecificProperties                                     string   `json:"less.lint.unknownVendorSpecificProperties,omitempty"`
	LessLintValidProperties                                                     []any    `json:"less.lint.validProperties,omitempty"`
	LessLintVendorPrefix                                                        string   `json:"less.lint.vendorPrefix,omitempty"`
	LessLintZeroUnits                                                           string   `json:"less.lint.zeroUnits,omitempty"`
	LessValidate                                                                bool     `json:"less.validate,omitempty"`
	ScssCompletionCompletePropertyWithSemicolon                                 bool     `json:"scss.completion.completePropertyWithSemicolon,omitempty"`
	ScssCompletionTriggerPropertyValueCompletion                                bool     `json:"scss.completion.triggerPropertyValueCompletion,omitempty"`
	ScssFormatBraceStyle                                                        string   `json:"scss.format.braceStyle,omitempty"`
	ScssFormatEnable                                                            bool     `json:"scss.format.enable,omitempty"`
	ScssFormatMaxPreserveNewLines                                               any      `json:"scss.format.maxPreserveNewLines,omitempty"`
	ScssFormatNewlineBetweenRules                                               bool     `json:"scss.format.newlineBetweenRules,omitempty"`
	ScssFormatNewlineBetweenSelectors                                           bool     `json:"scss.format.newlineBetweenSelectors,omitempty"`
	ScssFormatPreserveNewLines                                                  bool     `json:"scss.format.preserveNewLines,omitempty"`
	ScssFormatSpaceAroundSelectorSeparator                                      bool     `json:"scss.format.spaceAroundSelectorSeparator,omitempty"`
	ScssHoverDocumentation                                                      bool     `json:"scss.hover.documentation,omitempty"`
	ScssHoverReferences                                                         bool     `json:"scss.hover.references,omitempty"`
	ScssLintArgumentsInColorFunction                                            string   `json:"scss.lint.argumentsInColorFunction,omitempty"`
	ScssLintBoxModel                                                            string   `json:"scss.lint.boxModel,omitempty"`
	ScssLintCompatibleVendorPrefixes                                            string   `json:"scss.lint.compatibleVendorPrefixes,omitempty"`
	ScssLintDuplicateProperties                                                 string   `json:"scss.lint.duplicateProperties,omitempty"`
	ScssLintEmptyRules                                                          string   `json:"scss.lint.emptyRules,omitempty"`
	ScssLintFloat                                                               string   `json:"scss.lint.float,omitempty"`
	ScssLintFontFaceProperties                                                  string   `json:"scss.lint.fontFaceProperties,omitempty"`
	ScssLintHexColorLength                                                      string   `json:"scss.lint.hexColorLength,omitempty"`
	ScssLintIDSelector                                                          string   `json:"scss.lint.idSelector,omitempty"`
	ScssLintIeHack                                                              string   `json:"scss.lint.ieHack,omitempty"`
	ScssLintImportant                                                           string   `json:"scss.lint.important,omitempty"`
	ScssLintImportStatement                                                     string   `json:"scss.lint.importStatement,omitempty"`
	ScssLintPropertyIgnoredDueToDisplay                                         string   `json:"scss.lint.propertyIgnoredDueToDisplay,omitempty"`
	ScssLintUniversalSelector                                                   string   `json:"scss.lint.universalSelector,omitempty"`
	ScssLintUnknownAtRules                                                      string   `json:"scss.lint.unknownAtRules,omitempty"`
	ScssLintUnknownProperties                                                   string   `json:"scss.lint.unknownProperties,omitempty"`
	ScssLintUnknownVendorSpecificProperties                                     string   `json:"scss.lint.unknownVendorSpecificProperties,omitempty"`
	ScssLintValidProperties                                                     []any    `json:"scss.lint.validProperties,omitempty"`
	ScssLintVendorPrefix                                                        string   `json:"scss.lint.vendorPrefix,omitempty"`
	ScssLintZeroUnits                                                           string   `json:"scss.lint.zeroUnits,omitempty"`
	ScssValidate                                                                bool     `json:"scss.validate,omitempty"`
	ExtensionsAutoCheckUpdates                                                  bool     `json:"extensions.autoCheckUpdates,omitempty"`
	ExtensionsAutoUpdate                                                        bool     `json:"extensions.autoUpdate,omitempty"`
	ExtensionsCloseExtensionDetailsOnViewChange                                 bool     `json:"extensions.closeExtensionDetailsOnViewChange,omitempty"`
	ExtensionsConfirmedURIHandlerExtensionIds                                   []any    `json:"extensions.confirmedUriHandlerExtensionIds,omitempty"`
	ExtensionsIgnoreRecommendations                                             bool     `json:"extensions.ignoreRecommendations,omitempty"`
	ExtensionsSupportUntrustedWorkspaces                                        struct{} `json:"extensions.supportUntrustedWorkspaces,omitempty"`
	ExtensionsSupportVirtualWorkspaces                                          struct{} `json:"extensions.supportVirtualWorkspaces,omitempty"`
	OutputSmartScrollEnabled                                                    bool     `json:"output.smartScroll.enabled,omitempty"`
	SettingsSyncIgnoredExtensions                                               []any    `json:"settingsSync.ignoredExtensions,omitempty"`
	SettingsSyncIgnoredSettings                                                 []any    `json:"settingsSync.ignoredSettings,omitempty"`
	SettingsSyncKeybindingsPerPlatform                                          bool     `json:"settingsSync.keybindingsPerPlatform,omitempty"`
	InteractiveWindowCollapseCellInputCode                                      string   `json:"interactiveWindow.collapseCellInputCode,omitempty"`
	NotebookBreadcrumbsShowCodeCells                                            bool     `json:"notebook.breadcrumbs.showCodeCells,omitempty"`
	NotebookCellFocusIndicator                                                  string   `json:"notebook.cellFocusIndicator,omitempty"`
	NotebookCellToolbarLocation                                                 struct {
		Default string `json:"default,omitempty"`
	} `json:"notebook.cellToolbarLocation,omitempty"`
	NotebookCellToolbarVisibility       string   `json:"notebook.cellToolbarVisibility,omitempty"`
	NotebookCompactView                 bool     `json:"notebook.compactView,omitempty"`
	NotebookConfirmDeleteRunningCell    bool     `json:"notebook.confirmDeleteRunningCell,omitempty"`
	NotebookConsolidatedOutputButton    bool     `json:"notebook.consolidatedOutputButton,omitempty"`
	NotebookConsolidatedRunButton       bool     `json:"notebook.consolidatedRunButton,omitempty"`
	NotebookDiffEnablePreview           bool     `json:"notebook.diff.enablePreview,omitempty"`
	NotebookDiffIgnoreMetadata          bool     `json:"notebook.diff.ignoreMetadata,omitempty"`
	NotebookDiffIgnoreOutputs           bool     `json:"notebook.diff.ignoreOutputs,omitempty"`
	NotebookDiffOverviewRuler           bool     `json:"notebook.diff.overviewRuler,omitempty"`
	NotebookDisplayOrder                []any    `json:"notebook.displayOrder,omitempty"`
	NotebookDragAndDropEnabled          bool     `json:"notebook.dragAndDropEnabled,omitempty"`
	NotebookEditorOptionsCustomizations struct{} `json:"notebook.editorOptionsCustomizations,omitempty"`
	NotebookFindScope                   struct {
		MarkupSource  bool `json:"markupSource,omitempty"`
		MarkupPreview bool `json:"markupPreview,omitempty"`
		CodeSource    bool `json:"codeSource,omitempty"`
		CodeOutput    bool `json:"codeOutput,omitempty"`
	} `json:"notebook.find.scope,omitempty"`
	NotebookFormatOnCellExecution                        bool     `json:"notebook.formatOnCellExecution,omitempty"`
	NotebookFormatOnSaveEnabled                          bool     `json:"notebook.formatOnSave.enabled,omitempty"`
	NotebookGlobalToolbar                                bool     `json:"notebook.globalToolbar,omitempty"`
	NotebookGlobalToolbarShowLabel                       string   `json:"notebook.globalToolbarShowLabel,omitempty"`
	NotebookInsertToolbarLocation                        string   `json:"notebook.insertToolbarLocation,omitempty"`
	NotebookLineNumbers                                  string   `json:"notebook.lineNumbers,omitempty"`
	NotebookMarkupFontSize                               int      `json:"notebook.markup.fontSize,omitempty"`
	NotebookNavigationAllowNavigateToSurroundingCells    bool     `json:"notebook.navigation.allowNavigateToSurroundingCells,omitempty"`
	NotebookOutlineShowCodeCells                         bool     `json:"notebook.outline.showCodeCells,omitempty"`
	NotebookOutputFontFamily                             string   `json:"notebook.output.fontFamily,omitempty"`
	NotebookOutputFontSize                               int      `json:"notebook.output.fontSize,omitempty"`
	NotebookOutputLineHeight                             int      `json:"notebook.output.lineHeight,omitempty"`
	NotebookOutputScrolling                              bool     `json:"notebook.output.scrolling,omitempty"`
	NotebookOutputTextLineLimit                          int      `json:"notebook.output.textLineLimit,omitempty"`
	NotebookOutputWordWrap                               bool     `json:"notebook.output.wordWrap,omitempty"`
	NotebookShowCellStatusBar                            string   `json:"notebook.showCellStatusBar,omitempty"`
	NotebookShowFoldingControls                          string   `json:"notebook.showFoldingControls,omitempty"`
	NotebookUndoRedoPerCell                              bool     `json:"notebook.undoRedoPerCell,omitempty"`
	InteractiveWindowAlwaysScrollOnNewCell               bool     `json:"interactiveWindow.alwaysScrollOnNewCell,omitempty"`
	TerminalExplorerKind                                 string   `json:"terminal.explorerKind,omitempty"`
	TerminalExternalLinuxExec                            string   `json:"terminal.external.linuxExec,omitempty"`
	TerminalExternalOsxExec                              string   `json:"terminal.external.osxExec,omitempty"`
	TerminalExternalWindowsExec                          string   `json:"terminal.external.windowsExec,omitempty"`
	TerminalSourceControlRepositoriesKind                string   `json:"terminal.sourceControlRepositoriesKind,omitempty"`
	TerminalIntegratedAllowChords                        bool     `json:"terminal.integrated.allowChords,omitempty"`
	TerminalIntegratedAllowMnemonics                     bool     `json:"terminal.integrated.allowMnemonics,omitempty"`
	TerminalIntegratedAltClickMovesCursor                bool     `json:"terminal.integrated.altClickMovesCursor,omitempty"`
	TerminalIntegratedAutomationProfileLinux             any      `json:"terminal.integrated.automationProfile.linux,omitempty"`
	TerminalIntegratedAutomationProfileOsx               any      `json:"terminal.integrated.automationProfile.osx,omitempty"`
	TerminalIntegratedAutomationProfileWindows           any      `json:"terminal.integrated.automationProfile.windows,omitempty"`
	TerminalIntegratedAutoReplies                        struct{} `json:"terminal.integrated.autoReplies,omitempty"`
	TerminalIntegratedBellDuration                       int      `json:"terminal.integrated.bellDuration,omitempty"`
	TerminalIntegratedCommandsToSkipShell                []any    `json:"terminal.integrated.commandsToSkipShell,omitempty"`
	TerminalIntegratedConfirmOnExit                      string   `json:"terminal.integrated.confirmOnExit,omitempty"`
	TerminalIntegratedConfirmOnKill                      string   `json:"terminal.integrated.confirmOnKill,omitempty"`
	TerminalIntegratedCopyOnSelection                    bool     `json:"terminal.integrated.copyOnSelection,omitempty"`
	TerminalIntegratedCursorBlinking                     bool     `json:"terminal.integrated.cursorBlinking,omitempty"`
	TerminalIntegratedCursorStyle                        string   `json:"terminal.integrated.cursorStyle,omitempty"`
	TerminalIntegratedCursorWidth                        int      `json:"terminal.integrated.cursorWidth,omitempty"`
	TerminalIntegratedCustomGlyphs                       bool     `json:"terminal.integrated.customGlyphs,omitempty"`
	TerminalIntegratedCwd                                string   `json:"terminal.integrated.cwd,omitempty"`
	TerminalIntegratedDefaultLocation                    string   `json:"terminal.integrated.defaultLocation,omitempty"`
	TerminalIntegratedDefaultProfileLinux                any      `json:"terminal.integrated.defaultProfile.linux,omitempty"`
	TerminalIntegratedDefaultProfileOsx                  any      `json:"terminal.integrated.defaultProfile.osx,omitempty"`
	TerminalIntegratedDefaultProfileWindows              any      `json:"terminal.integrated.defaultProfile.windows,omitempty"`
	TerminalIntegratedDetectLocale                       string   `json:"terminal.integrated.detectLocale,omitempty"`
	TerminalIntegratedDrawBoldTextInBrightColors         bool     `json:"terminal.integrated.drawBoldTextInBrightColors,omitempty"`
	TerminalIntegratedEnableBell                         bool     `json:"terminal.integrated.enableBell,omitempty"`
	TerminalIntegratedEnableFileLinks                    string   `json:"terminal.integrated.enableFileLinks,omitempty"`
	TerminalIntegratedEnableImages                       bool     `json:"terminal.integrated.enableImages,omitempty"`
	TerminalIntegratedEnableMultiLinePasteWarning        bool     `json:"terminal.integrated.enableMultiLinePasteWarning,omitempty"`
	TerminalIntegratedEnablePersistentSessions           bool     `json:"terminal.integrated.enablePersistentSessions,omitempty"`
	TerminalIntegratedEnvLinux                           struct{} `json:"terminal.integrated.env.linux,omitempty"`
	TerminalIntegratedEnvOsx                             struct{} `json:"terminal.integrated.env.osx,omitempty"`
	TerminalIntegratedEnvWindows                         struct{} `json:"terminal.integrated.env.windows,omitempty"`
	TerminalIntegratedEnvironmentChangesIndicator        string   `json:"terminal.integrated.environmentChangesIndicator,omitempty"`
	TerminalIntegratedEnvironmentChangesRelaunch         bool     `json:"terminal.integrated.environmentChangesRelaunch,omitempty"`
	TerminalIntegratedFastScrollSensitivity              int      `json:"terminal.integrated.fastScrollSensitivity,omitempty"`
	TerminalIntegratedFontFamily                         string   `json:"terminal.integrated.fontFamily,omitempty"`
	TerminalIntegratedFontSize                           int      `json:"terminal.integrated.fontSize,omitempty"`
	TerminalIntegratedFontWeight                         string   `json:"terminal.integrated.fontWeight,omitempty"`
	TerminalIntegratedFontWeightBold                     string   `json:"terminal.integrated.fontWeightBold,omitempty"`
	TerminalIntegratedGpuAcceleration                    string   `json:"terminal.integrated.gpuAcceleration,omitempty"`
	TerminalIntegratedIgnoreProcessNames                 []any    `json:"terminal.integrated.ignoreProcessNames,omitempty"`
	TerminalIntegratedInheritEnv                         bool     `json:"terminal.integrated.inheritEnv,omitempty"`
	TerminalIntegratedLetterSpacing                      int      `json:"terminal.integrated.letterSpacing,omitempty"`
	TerminalIntegratedLineHeight                         int      `json:"terminal.integrated.lineHeight,omitempty"`
	TerminalIntegratedLocalEchoEnabled                   string   `json:"terminal.integrated.localEchoEnabled,omitempty"`
	TerminalIntegratedLocalEchoExcludePrograms           []string `json:"terminal.integrated.localEchoExcludePrograms,omitempty"`
	TerminalIntegratedLocalEchoLatencyThreshold          int      `json:"terminal.integrated.localEchoLatencyThreshold,omitempty"`
	TerminalIntegratedLocalEchoStyle                     string   `json:"terminal.integrated.localEchoStyle,omitempty"`
	TerminalIntegratedMacOptionClickForcesSelection      bool     `json:"terminal.integrated.macOptionClickForcesSelection,omitempty"`
	TerminalIntegratedMacOptionIsMeta                    bool     `json:"terminal.integrated.macOptionIsMeta,omitempty"`
	TerminalIntegratedMinimumContrastRatio               float64  `json:"terminal.integrated.minimumContrastRatio,omitempty"`
	TerminalIntegratedMouseWheelScrollSensitivity        int      `json:"terminal.integrated.mouseWheelScrollSensitivity,omitempty"`
	TerminalIntegratedPersistentSessionReviveProcess     string   `json:"terminal.integrated.persistentSessionReviveProcess,omitempty"`
	TerminalIntegratedPersistentSessionScrollback        int      `json:"terminal.integrated.persistentSessionScrollback,omitempty"`
	TerminalIntegratedProfilesLinux                      struct{} `json:"terminal.integrated.profiles.linux,omitempty"`
	TerminalIntegratedProfilesOsx                        struct{} `json:"terminal.integrated.profiles.osx,omitempty"`
	TerminalIntegratedProfilesWindows                    struct{} `json:"terminal.integrated.profiles.windows,omitempty"`
	TerminalIntegratedRightClickBehavior                 string   `json:"terminal.integrated.rightClickBehavior,omitempty"`
	TerminalIntegratedScrollback                         int      `json:"terminal.integrated.scrollback,omitempty"`
	TerminalIntegratedSendKeybindingsToShell             bool     `json:"terminal.integrated.sendKeybindingsToShell,omitempty"`
	TerminalIntegratedShellIntegrationDecorationsEnabled string   `json:"terminal.integrated.shellIntegration.decorationsEnabled,omitempty"`
	TerminalIntegratedShellIntegrationEnabled            bool     `json:"terminal.integrated.shellIntegration.enabled,omitempty"`
	TerminalIntegratedShellIntegrationHistory            int      `json:"terminal.integrated.shellIntegration.history,omitempty"`
	TerminalIntegratedShowExitAlert                      bool     `json:"terminal.integrated.showExitAlert,omitempty"`
	TerminalIntegratedShowLinkHover                      bool     `json:"terminal.integrated.showLinkHover,omitempty"`
	TerminalIntegratedSmoothScrolling                    bool     `json:"terminal.integrated.smoothScrolling,omitempty"`
	TerminalIntegratedSplitCwd                           string   `json:"terminal.integrated.splitCwd,omitempty"`
	TerminalIntegratedTabFocusMode                       any      `json:"terminal.integrated.tabFocusMode,omitempty"`
	TerminalIntegratedTabsDefaultColor                   any      `json:"terminal.integrated.tabs.defaultColor,omitempty"`
	TerminalIntegratedTabsDefaultIcon                    string   `json:"terminal.integrated.tabs.defaultIcon,omitempty"`
	TerminalIntegratedTabsDescription                    string   `json:"terminal.integrated.tabs.description,omitempty"`
	TerminalIntegratedTabsEnableAnimation                bool     `json:"terminal.integrated.tabs.enableAnimation,omitempty"`
	TerminalIntegratedTabsEnabled                        bool     `json:"terminal.integrated.tabs.enabled,omitempty"`
	TerminalIntegratedTabsFocusMode                      string   `json:"terminal.integrated.tabs.focusMode,omitempty"`
	TerminalIntegratedTabsHideCondition                  string   `json:"terminal.integrated.tabs.hideCondition,omitempty"`
	TerminalIntegratedTabsLocation                       string   `json:"terminal.integrated.tabs.location,omitempty"`
	TerminalIntegratedTabsSeparator                      string   `json:"terminal.integrated.tabs.separator,omitempty"`
	TerminalIntegratedTabsShowActions                    string   `json:"terminal.integrated.tabs.showActions,omitempty"`
	TerminalIntegratedTabsShowActiveTerminal             string   `json:"terminal.integrated.tabs.showActiveTerminal,omitempty"`
	TerminalIntegratedTabsTitle                          string   `json:"terminal.integrated.tabs.title,omitempty"`
	TerminalIntegratedTabStopWidth                       int      `json:"terminal.integrated.tabStopWidth,omitempty"`
	TerminalIntegratedUnicodeVersion                     string   `json:"terminal.integrated.unicodeVersion,omitempty"`
	TerminalIntegratedUseWslProfiles                     bool     `json:"terminal.integrated.useWslProfiles,omitempty"`
	TerminalIntegratedWindowsEnableConpty                bool     `json:"terminal.integrated.windowsEnableConpty,omitempty"`
	TerminalIntegratedWordSeparators                     string   `json:"terminal.integrated.wordSeparators,omitempty"`
	TaskAllowAutomaticTasks                              string   `json:"task.allowAutomaticTasks,omitempty"`
	TaskAutoDetect                                       string   `json:"task.autoDetect,omitempty"`
	TaskProblemMatchersNeverPrompt                       bool     `json:"task.problemMatchers.neverPrompt,omitempty"`
	TaskQuickOpenDetail                                  bool     `json:"task.quickOpen.detail,omitempty"`
	TaskQuickOpenHistory                                 int      `json:"task.quickOpen.history,omitempty"`
	TaskQuickOpenShowAll                                 bool     `json:"task.quickOpen.showAll,omitempty"`
	TaskQuickOpenSkip                                    bool     `json:"task.quickOpen.skip,omitempty"`
	TaskReconnection                                     bool     `json:"task.reconnection,omitempty"`
	TaskSaveBeforeRun                                    string   `json:"task.saveBeforeRun,omitempty"`
	TaskSlowProviderWarning                              bool     `json:"task.slowProviderWarning,omitempty"`
	ProblemsAutoReveal                                   bool     `json:"problems.autoReveal,omitempty"`
	ProblemsDecorationsEnabled                           bool     `json:"problems.decorations.enabled,omitempty"`
	ProblemsDefaultViewMode                              string   `json:"problems.defaultViewMode,omitempty"`
	ProblemsShowCurrentInStatus                          bool     `json:"problems.showCurrentInStatus,omitempty"`
	ProblemsSortOrder                                    string   `json:"problems.sortOrder,omitempty"`
	BreadcrumbsEnabled                                   bool     `json:"breadcrumbs.enabled,omitempty"`
	BreadcrumbsFilePath                                  string   `json:"breadcrumbs.filePath,omitempty"`
	BreadcrumbsIcons                                     bool     `json:"breadcrumbs.icons,omitempty"`
	BreadcrumbsShowArrays                                bool     `json:"breadcrumbs.showArrays,omitempty"`
	BreadcrumbsShowBooleans                              bool     `json:"breadcrumbs.showBooleans,omitempty"`
	BreadcrumbsShowClasses                               bool     `json:"breadcrumbs.showClasses,omitempty"`
	BreadcrumbsShowConstants                             bool     `json:"breadcrumbs.showConstants,omitempty"`
	BreadcrumbsShowConstructors                          bool     `json:"breadcrumbs.showConstructors,omitempty"`
	BreadcrumbsShowEnumMembers                           bool     `json:"breadcrumbs.showEnumMembers,omitempty"`
	BreadcrumbsShowEnums                                 bool     `json:"breadcrumbs.showEnums,omitempty"`
	BreadcrumbsShowEvents                                bool     `json:"breadcrumbs.showEvents,omitempty"`
	BreadcrumbsShowFields                                bool     `json:"breadcrumbs.showFields,omitempty"`
	BreadcrumbsShowFiles                                 bool     `json:"breadcrumbs.showFiles,omitempty"`
	BreadcrumbsShowFunctions                             bool     `json:"breadcrumbs.showFunctions,omitempty"`
	BreadcrumbsShowInterfaces                            bool     `json:"breadcrumbs.showInterfaces,omitempty"`
	BreadcrumbsShowKeys                                  bool     `json:"breadcrumbs.showKeys,omitempty"`
	BreadcrumbsShowMethods                               bool     `json:"breadcrumbs.showMethods,omitempty"`
	BreadcrumbsShowModules                               bool     `json:"breadcrumbs.showModules,omitempty"`
	BreadcrumbsShowNamespaces                            bool     `json:"breadcrumbs.showNamespaces,omitempty"`
	BreadcrumbsShowNull                                  bool     `json:"breadcrumbs.showNull,omitempty"`
	BreadcrumbsShowNumbers                               bool     `json:"breadcrumbs.showNumbers,omitempty"`
	BreadcrumbsShowObjects                               bool     `json:"breadcrumbs.showObjects,omitempty"`
	BreadcrumbsShowOperators                             bool     `json:"breadcrumbs.showOperators,omitempty"`
	BreadcrumbsShowPackages                              bool     `json:"breadcrumbs.showPackages,omitempty"`
	BreadcrumbsShowProperties                            bool     `json:"breadcrumbs.showProperties,omitempty"`
	BreadcrumbsShowStrings                               bool     `json:"breadcrumbs.showStrings,omitempty"`
	BreadcrumbsShowStructs                               bool     `json:"breadcrumbs.showStructs,omitempty"`
	BreadcrumbsShowTypeParameters                        bool     `json:"breadcrumbs.showTypeParameters,omitempty"`
	BreadcrumbsShowVariables                             bool     `json:"breadcrumbs.showVariables,omitempty"`
	BreadcrumbsSymbolPath                                string   `json:"breadcrumbs.symbolPath,omitempty"`
	BreadcrumbsSymbolSortOrder                           string   `json:"breadcrumbs.symbolSortOrder,omitempty"`
	OutlineCollapseItems                                 string   `json:"outline.collapseItems,omitempty"`
	OutlineIcons                                         bool     `json:"outline.icons,omitempty"`
	OutlineProblemsBadges                                bool     `json:"outline.problems.badges,omitempty"`
	OutlineProblemsColors                                bool     `json:"outline.problems.colors,omitempty"`
	OutlineProblemsEnabled                               bool     `json:"outline.problems.enabled,omitempty"`
	OutlineShowArrays                                    bool     `json:"outline.showArrays,omitempty"`
	OutlineShowBooleans                                  bool     `json:"outline.showBooleans,omitempty"`
	OutlineShowClasses                                   bool     `json:"outline.showClasses,omitempty"`
	OutlineShowConstants                                 bool     `json:"outline.showConstants,omitempty"`
	OutlineShowConstructors                              bool     `json:"outline.showConstructors,omitempty"`
	OutlineShowEnumMembers                               bool     `json:"outline.showEnumMembers,omitempty"`
	OutlineShowEnums                                     bool     `json:"outline.showEnums,omitempty"`
	OutlineShowEvents                                    bool     `json:"outline.showEvents,omitempty"`
	OutlineShowFields                                    bool     `json:"outline.showFields,omitempty"`
	OutlineShowFiles                                     bool     `json:"outline.showFiles,omitempty"`
	OutlineShowFunctions                                 bool     `json:"outline.showFunctions,omitempty"`
	OutlineShowInterfaces                                bool     `json:"outline.showInterfaces,omitempty"`
	OutlineShowKeys                                      bool     `json:"outline.showKeys,omitempty"`
	OutlineShowMethods                                   bool     `json:"outline.showMethods,omitempty"`
	OutlineShowModules                                   bool     `json:"outline.showModules,omitempty"`
	OutlineShowNamespaces                                bool     `json:"outline.showNamespaces,omitempty"`
	OutlineShowNull                                      bool     `json:"outline.showNull,omitempty"`
	OutlineShowNumbers                                   bool     `json:"outline.showNumbers,omitempty"`
	OutlineShowObjects                                   bool     `json:"outline.showObjects,omitempty"`
	OutlineShowOperators                                 bool     `json:"outline.showOperators,omitempty"`
	OutlineShowPackages                                  bool     `json:"outline.showPackages,omitempty"`
	OutlineShowProperties                                bool     `json:"outline.showProperties,omitempty"`
	OutlineShowStrings                                   bool     `json:"outline.showStrings,omitempty"`
	OutlineShowStructs                                   bool     `json:"outline.showStructs,omitempty"`
	OutlineShowTypeParameters                            bool     `json:"outline.showTypeParameters,omitempty"`
	OutlineShowVariables                                 bool     `json:"outline.showVariables,omitempty"`
	TimelinePageSize                                     any      `json:"timeline.pageSize,omitempty"`
	Clojure                                              struct {
		DiffEditorIgnoreTrimWhitespace bool `json:"diffEditor.ignoreTrimWhitespace,omitempty"`
	} `json:"[clojure],omitempty"`
	Coffeescript struct {
		DiffEditorIgnoreTrimWhitespace bool `json:"diffEditor.ignoreTrimWhitespace,omitempty"`
	} `json:"[coffeescript],omitempty"`
	Csharp struct {
		EditorMaxTokenizationLineLength int `json:"editor.maxTokenizationLineLength,omitempty"`
	} `json:"[csharp],omitempty"`
	CSS struct {
		EditorSuggestInsertMode string `json:"editor.suggest.insertMode,omitempty"`
	} `json:"[css],omitempty"`
	Dockercompose struct {
		EditorInsertSpaces bool   `json:"editor.insertSpaces,omitempty"`
		EditorTabSize      int    `json:"editor.tabSize,omitempty"`
		EditorAutoIndent   string `json:"editor.autoIndent,omitempty"`
	} `json:"[dockercompose],omitempty"`
	Dockerfile struct {
		EditorQuickSuggestions struct {
			Strings bool `json:"strings,omitempty"`
		} `json:"editor.quickSuggestions,omitempty"`
	} `json:"[dockerfile],omitempty"`
	Fsharp struct {
		DiffEditorIgnoreTrimWhitespace bool `json:"diffEditor.ignoreTrimWhitespace,omitempty"`
	} `json:"[fsharp],omitempty"`
	GitCommit struct {
		EditorRulers                    []int `json:"editor.rulers,omitempty"`
		WorkbenchEditorRestoreViewState bool  `json:"workbench.editor.restoreViewState,omitempty"`
	} `json:"[git-commit],omitempty"`
	GitRebase struct {
		WorkbenchEditorRestoreViewState bool `json:"workbench.editor.restoreViewState,omitempty"`
	} `json:"[git-rebase],omitempty"`
	Go struct {
		EditorInsertSpaces bool `json:"editor.insertSpaces,omitempty"`
	} `json:"[go],omitempty"`
	Handlebars struct {
		EditorSuggestInsertMode string `json:"editor.suggest.insertMode,omitempty"`
	} `json:"[handlebars],omitempty"`
	HTML struct {
		EditorSuggestInsertMode string `json:"editor.suggest.insertMode,omitempty"`
	} `json:"[html],omitempty"`
	Jade struct {
		DiffEditorIgnoreTrimWhitespace bool `json:"diffEditor.ignoreTrimWhitespace,omitempty"`
	} `json:"[jade],omitempty"`
	Javascript struct {
		EditorMaxTokenizationLineLength int `json:"editor.maxTokenizationLineLength,omitempty"`
	} `json:"[javascript],omitempty"`
	JSON struct {
		EditorQuickSuggestions struct {
			Strings bool `json:"strings,omitempty"`
		} `json:"editor.quickSuggestions,omitempty"`
		EditorSuggestInsertMode string `json:"editor.suggest.insertMode,omitempty"`
	} `json:"[json],omitempty"`
	Jsonc struct {
		EditorQuickSuggestions struct {
			Strings bool `json:"strings,omitempty"`
		} `json:"editor.quickSuggestions,omitempty"`
		EditorSuggestInsertMode string `json:"editor.suggest.insertMode,omitempty"`
	} `json:"[jsonc],omitempty"`
	Less struct {
		EditorSuggestInsertMode string `json:"editor.suggest.insertMode,omitempty"`
	} `json:"[less],omitempty"`
	Makefile struct {
		EditorInsertSpaces bool `json:"editor.insertSpaces,omitempty"`
	} `json:"[makefile],omitempty"`
	Markdown struct {
		EditorUnicodeHighlightAmbiguousCharacters bool   `json:"editor.unicodeHighlight.ambiguousCharacters,omitempty"`
		EditorUnicodeHighlightInvisibleCharacters bool   `json:"editor.unicodeHighlight.invisibleCharacters,omitempty"`
		DiffEditorIgnoreTrimWhitespace            bool   `json:"diffEditor.ignoreTrimWhitespace,omitempty"`
		EditorWordWrap                            string `json:"editor.wordWrap,omitempty"`
		EditorQuickSuggestions                    struct {
			Comments string `json:"comments,omitempty"`
			Strings  string `json:"strings,omitempty"`
			Other    string `json:"other,omitempty"`
		} `json:"editor.quickSuggestions,omitempty"`
	} `json:"[markdown],omitempty"`
	Plaintext struct {
		EditorUnicodeHighlightAmbiguousCharacters bool `json:"editor.unicodeHighlight.ambiguousCharacters,omitempty"`
		EditorUnicodeHighlightInvisibleCharacters bool `json:"editor.unicodeHighlight.invisibleCharacters,omitempty"`
	} `json:"[plaintext],omitempty"`
	Python struct {
		DiffEditorIgnoreTrimWhitespace bool `json:"diffEditor.ignoreTrimWhitespace,omitempty"`
	} `json:"[python],omitempty"`
	Scss struct {
		EditorSuggestInsertMode string `json:"editor.suggest.insertMode,omitempty"`
	} `json:"[scss],omitempty"`
	SearchResult struct {
		EditorLineNumbers string `json:"editor.lineNumbers,omitempty"`
	} `json:"[search-result],omitempty"`
	Shellscript struct {
		FilesEol string `json:"files.eol,omitempty"`
	} `json:"[shellscript],omitempty"`
	Yaml struct {
		EditorInsertSpaces             bool   `json:"editor.insertSpaces,omitempty"`
		EditorTabSize                  int    `json:"editor.tabSize,omitempty"`
		EditorAutoIndent               string `json:"editor.autoIndent,omitempty"`
		DiffEditorIgnoreTrimWhitespace bool   `json:"diffEditor.ignoreTrimWhitespace,omitempty"`
	} `json:"[yaml],omitempty"`
	AudioCuesChatRequestSent            string `json:"audioCues.chatRequestSent,omitempty"`
	AudioCuesChatResponsePending        string `json:"audioCues.chatResponsePending,omitempty"`
	AudioCuesChatResponseReceived       string `json:"audioCues.chatResponseReceived,omitempty"`
	AudioCuesDebouncePositionChanges    bool   `json:"audioCues.debouncePositionChanges,omitempty"`
	AudioCuesDiffLineDeleted            string `json:"audioCues.diffLineDeleted,omitempty"`
	AudioCuesDiffLineInserted           string `json:"audioCues.diffLineInserted,omitempty"`
	AudioCuesDiffLineModified           string `json:"audioCues.diffLineModified,omitempty"`
	AudioCuesLineHasBreakpoint          string `json:"audioCues.lineHasBreakpoint,omitempty"`
	AudioCuesLineHasError               string `json:"audioCues.lineHasError,omitempty"`
	AudioCuesLineHasFoldedArea          string `json:"audioCues.lineHasFoldedArea,omitempty"`
	AudioCuesLineHasInlineSuggestion    string `json:"audioCues.lineHasInlineSuggestion,omitempty"`
	AudioCuesLineHasWarning             string `json:"audioCues.lineHasWarning,omitempty"`
	AudioCuesNoInlayHints               string `json:"audioCues.noInlayHints,omitempty"`
	AudioCuesNotebookCellCompleted      string `json:"audioCues.notebookCellCompleted,omitempty"`
	AudioCuesNotebookCellFailed         string `json:"audioCues.notebookCellFailed,omitempty"`
	AudioCuesOnDebugBreak               string `json:"audioCues.onDebugBreak,omitempty"`
	AudioCuesTaskCompleted              string `json:"audioCues.taskCompleted,omitempty"`
	AudioCuesTaskFailed                 string `json:"audioCues.taskFailed,omitempty"`
	AudioCuesTerminalCommandFailed      string `json:"audioCues.terminalCommandFailed,omitempty"`
	AudioCuesTerminalQuickFix           string `json:"audioCues.terminalQuickFix,omitempty"`
	AudioCuesVolume                     int    `json:"audioCues.volume,omitempty"`
	RemoteTunnelsAccessHostNameOverride string `json:"remote.tunnels.access.hostNameOverride,omitempty"`
	RemoteTunnelsAccessPreventSleep     bool   `json:"remote.tunnels.access.preventSleep,omitempty"`
	RemoteAutoForwardPorts              bool   `json:"remote.autoForwardPorts,omitempty"`
	RemoteAutoForwardPortsSource        string `json:"remote.autoForwardPortsSource,omitempty"`
	RemoteDownloadExtensionsLocally     bool   `json:"remote.downloadExtensionsLocally,omitempty"`
	RemoteExtensionKind                 struct {
		PubName []string `json:"pub.name,omitempty"`
	} `json:"remote.extensionKind,omitempty"`
	RemoteForwardOnOpen                     bool     `json:"remote.forwardOnOpen,omitempty"`
	RemoteLocalPortHost                     string   `json:"remote.localPortHost,omitempty"`
	RemoteOtherPortsAttributes              struct{} `json:"remote.otherPortsAttributes,omitempty"`
	RemotePortsAttributes                   struct{} `json:"remote.portsAttributes,omitempty"`
	RemoteRestoreForwardedPorts             bool     `json:"remote.restoreForwardedPorts,omitempty"`
	AccessibilityVerbosityDiffEditor        bool     `json:"accessibility.verbosity.diffEditor,omitempty"`
	AccessibilityVerbosityHover             bool     `json:"accessibility.verbosity.hover,omitempty"`
	AccessibilityVerbosityInlineChat        bool     `json:"accessibility.verbosity.inlineChat,omitempty"`
	AccessibilityVerbosityKeybindingsEditor bool     `json:"accessibility.verbosity.keybindingsEditor,omitempty"`
	AccessibilityVerbosityNotebook          bool     `json:"accessibility.verbosity.notebook,omitempty"`
	AccessibilityVerbosityNotification      bool     `json:"accessibility.verbosity.notification,omitempty"`
	AccessibilityVerbosityPanelChat         bool     `json:"accessibility.verbosity.panelChat,omitempty"`
	AccessibilityVerbosityTerminal          bool     `json:"accessibility.verbosity.terminal,omitempty"`
	MergeEditorDiffAlgorithm                string   `json:"mergeEditor.diffAlgorithm,omitempty"`
	MergeEditorShowDeletionMarkers          bool     `json:"mergeEditor.showDeletionMarkers,omitempty"`
	EmmetExcludeLanguages                   []string `json:"emmet.excludeLanguages,omitempty"`
	EmmetExtensionsPath                     []any    `json:"emmet.extensionsPath,omitempty"`
	EmmetIncludeLanguages                   struct{} `json:"emmet.includeLanguages,omitempty"`
	EmmetOptimizeStylesheetParsing          bool     `json:"emmet.optimizeStylesheetParsing,omitempty"`
	EmmetPreferences                        struct{} `json:"emmet.preferences,omitempty"`
	EmmetShowAbbreviationSuggestions        bool     `json:"emmet.showAbbreviationSuggestions,omitempty"`
	EmmetShowExpandedAbbreviation           string   `json:"emmet.showExpandedAbbreviation,omitempty"`
	EmmetShowSuggestionsAsSnippets          bool     `json:"emmet.showSuggestionsAsSnippets,omitempty"`
	EmmetSyntaxProfiles                     struct{} `json:"emmet.syntaxProfiles,omitempty"`
	EmmetTriggerExpansionOnTab              bool     `json:"emmet.triggerExpansionOnTab,omitempty"`
	EmmetUseInlineCompletions               bool     `json:"emmet.useInlineCompletions,omitempty"`
	EmmetVariables                          struct{} `json:"emmet.variables,omitempty"`
	GitAllowForcePush                       bool     `json:"git.allowForcePush,omitempty"`
	GitAllowNoVerifyCommit                  bool     `json:"git.allowNoVerifyCommit,omitempty"`
	GitAlwaysShowStagedChangesResourceGroup bool     `json:"git.alwaysShowStagedChangesResourceGroup,omitempty"`
	GitAlwaysSignOff                        bool     `json:"git.alwaysSignOff,omitempty"`
	GitAutofetch                            bool     `json:"git.autofetch,omitempty"`
	GitAutofetchPeriod                      int      `json:"git.autofetchPeriod,omitempty"`
	GitAutorefresh                          bool     `json:"git.autorefresh,omitempty"`
	GitAutoRepositoryDetection              bool     `json:"git.autoRepositoryDetection,omitempty"`
	GitAutoStash                            bool     `json:"git.autoStash,omitempty"`
	GitBranchPrefix                         string   `json:"git.branchPrefix,omitempty"`
	GitBranchProtection                     []any    `json:"git.branchProtection,omitempty"`
	GitBranchProtectionPrompt               string   `json:"git.branchProtectionPrompt,omitempty"`
	GitBranchRandomNameDictionary           []string `json:"git.branchRandomName.dictionary,omitempty"`
	GitBranchRandomNameEnable               bool     `json:"git.branchRandomName.enable,omitempty"`
	GitBranchSortOrder                      string   `json:"git.branchSortOrder,omitempty"`
	GitBranchValidationRegex                string   `json:"git.branchValidationRegex,omitempty"`
	GitBranchWhitespaceChar                 string   `json:"git.branchWhitespaceChar,omitempty"`
	GitCheckoutType                         []string `json:"git.checkoutType,omitempty"`
	GitCloseDiffOnOperation                 bool     `json:"git.closeDiffOnOperation,omitempty"`
	GitCommandsToLog                        []any    `json:"git.commandsToLog,omitempty"`
	GitConfirmEmptyCommits                  bool     `json:"git.confirmEmptyCommits,omitempty"`
	GitConfirmForcePush                     bool     `json:"git.confirmForcePush,omitempty"`
	GitConfirmNoVerifyCommit                bool     `json:"git.confirmNoVerifyCommit,omitempty"`
	GitConfirmSync                          bool     `json:"git.confirmSync,omitempty"`
	GitCountBadge                           string   `json:"git.countBadge,omitempty"`
	GitDecorationsEnabled                   bool     `json:"git.decorations.enabled,omitempty"`
	GitDefaultBranchName                    string   `json:"git.defaultBranchName,omitempty"`
	GitDefaultCloneDirectory                any      `json:"git.defaultCloneDirectory,omitempty"`
	GitDetectSubmodules                     bool     `json:"git.detectSubmodules,omitempty"`
	GitDetectSubmodulesLimit                int      `json:"git.detectSubmodulesLimit,omitempty"`
	GitEnableCommitSigning                  bool     `json:"git.enableCommitSigning,omitempty"`
	GitEnabled                              bool     `json:"git.enabled,omitempty"`
	GitEnableSmartCommit                    bool     `json:"git.enableSmartCommit,omitempty"`
	GitEnableStatusBarSync                  bool     `json:"git.enableStatusBarSync,omitempty"`
	GitFetchOnPull                          bool     `json:"git.fetchOnPull,omitempty"`
	GitFollowTagsWhenSync                   bool     `json:"git.followTagsWhenSync,omitempty"`
	GitIgnoredRepositories                  []any    `json:"git.ignoredRepositories,omitempty"`
	GitIgnoreLegacyWarning                  bool     `json:"git.ignoreLegacyWarning,omitempty"`
	GitIgnoreLimitWarning                   bool     `json:"git.ignoreLimitWarning,omitempty"`
	GitIgnoreMissingGitWarning              bool     `json:"git.ignoreMissingGitWarning,omitempty"`
	GitIgnoreRebaseWarning                  bool     `json:"git.ignoreRebaseWarning,omitempty"`
	GitIgnoreSubmodules                     bool     `json:"git.ignoreSubmodules,omitempty"`
	GitIgnoreWindowsGit27Warning            bool     `json:"git.ignoreWindowsGit27Warning,omitempty"`
	GitInputValidation                      string   `json:"git.inputValidation,omitempty"`
	GitInputValidationLength                int      `json:"git.inputValidationLength,omitempty"`
	GitInputValidationSubjectLength         int      `json:"git.inputValidationSubjectLength,omitempty"`
	GitMergeEditor                          bool     `json:"git.mergeEditor,omitempty"`
	GitOpenAfterClone                       string   `json:"git.openAfterClone,omitempty"`
	GitOpenDiffOnClick                      bool     `json:"git.openDiffOnClick,omitempty"`
	GitOpenRepositoryInParentFolders        string   `json:"git.openRepositoryInParentFolders,omitempty"`
	GitOptimisticUpdate                     bool     `json:"git.optimisticUpdate,omitempty"`
	GitPath                                 any      `json:"git.path,omitempty"`
	GitPostCommitCommand                    string   `json:"git.postCommitCommand,omitempty"`
	GitPromptToSaveFilesBeforeCommit        string   `json:"git.promptToSaveFilesBeforeCommit,omitempty"`
	GitPromptToSaveFilesBeforeStash         string   `json:"git.promptToSaveFilesBeforeStash,omitempty"`
	GitPruneOnFetch                         bool     `json:"git.pruneOnFetch,omitempty"`
	GitPullBeforeCheckout                   bool     `json:"git.pullBeforeCheckout,omitempty"`
	GitPullTags                             bool     `json:"git.pullTags,omitempty"`
	GitRebaseWhenSync                       bool     `json:"git.rebaseWhenSync,omitempty"`
	GitRememberPostCommitCommand            bool     `json:"git.rememberPostCommitCommand,omitempty"`
	GitRepositoryScanIgnoredFolders         []string `json:"git.repositoryScanIgnoredFolders,omitempty"`
	GitRepositoryScanMaxDepth               int      `json:"git.repositoryScanMaxDepth,omitempty"`
	GitRequireGitUserConfig                 bool     `json:"git.requireGitUserConfig,omitempty"`
	GitScanRepositories                     []any    `json:"git.scanRepositories,omitempty"`
	GitShowActionButton                     struct {
		Commit  bool `json:"commit,omitempty"`
		Publish bool `json:"publish,omitempty"`
		Sync    bool `json:"sync,omitempty"`
	} `json:"git.showActionButton,omitempty"`
	GitShowCommitInput                             bool     `json:"git.showCommitInput,omitempty"`
	GitShowInlineOpenFileAction                    bool     `json:"git.showInlineOpenFileAction,omitempty"`
	GitShowProgress                                bool     `json:"git.showProgress,omitempty"`
	GitShowPushSuccessNotification                 bool     `json:"git.showPushSuccessNotification,omitempty"`
	GitSimilarityThreshold                         int      `json:"git.similarityThreshold,omitempty"`
	GitSmartCommitChanges                          string   `json:"git.smartCommitChanges,omitempty"`
	GitStatusLimit                                 int      `json:"git.statusLimit,omitempty"`
	GitSuggestSmartCommit                          bool     `json:"git.suggestSmartCommit,omitempty"`
	GitSupportCancellation                         bool     `json:"git.supportCancellation,omitempty"`
	GitTerminalAuthentication                      bool     `json:"git.terminalAuthentication,omitempty"`
	GitTerminalGitEditor                           bool     `json:"git.terminalGitEditor,omitempty"`
	GitTimelineDate                                string   `json:"git.timeline.date,omitempty"`
	GitTimelineShowAuthor                          bool     `json:"git.timeline.showAuthor,omitempty"`
	GitTimelineShowUncommitted                     bool     `json:"git.timeline.showUncommitted,omitempty"`
	GitUntrackedChanges                            string   `json:"git.untrackedChanges,omitempty"`
	GitUseCommitInputAsStashMessage                bool     `json:"git.useCommitInputAsStashMessage,omitempty"`
	GitUseEditorAsCommitInput                      bool     `json:"git.useEditorAsCommitInput,omitempty"`
	GitUseForcePushWithLease                       bool     `json:"git.useForcePushWithLease,omitempty"`
	GitUseIntegratedAskPass                        bool     `json:"git.useIntegratedAskPass,omitempty"`
	GitVerboseCommit                               bool     `json:"git.verboseCommit,omitempty"`
	GithubBranchProtection                         bool     `json:"github.branchProtection,omitempty"`
	GithubGitAuthentication                        bool     `json:"github.gitAuthentication,omitempty"`
	GithubGitProtocol                              string   `json:"github.gitProtocol,omitempty"`
	GithubEnterpriseURI                            string   `json:"github-enterprise.uri,omitempty"`
	GruntAutoDetect                                string   `json:"grunt.autoDetect,omitempty"`
	GulpAutoDetect                                 string   `json:"gulp.autoDetect,omitempty"`
	JakeAutoDetect                                 string   `json:"jake.autoDetect,omitempty"`
	MediaPreviewVideoAutoPlay                      bool     `json:"mediaPreview.video.autoPlay,omitempty"`
	MediaPreviewVideoLoop                          bool     `json:"mediaPreview.video.loop,omitempty"`
	MergeConflictAutoNavigateNextConflictEnabled   bool     `json:"merge-conflict.autoNavigateNextConflict.enabled,omitempty"`
	MergeConflictCodeLensEnabled                   bool     `json:"merge-conflict.codeLens.enabled,omitempty"`
	MergeConflictDecoratorsEnabled                 bool     `json:"merge-conflict.decorators.enabled,omitempty"`
	MergeConflictDiffViewPosition                  string   `json:"merge-conflict.diffViewPosition,omitempty"`
	MicrosoftSovereignCloudCustomEnvironment       struct{} `json:"microsoft-sovereign-cloud.customEnvironment,omitempty"`
	MicrosoftSovereignCloudEnvironment             string   `json:"microsoft-sovereign-cloud.environment,omitempty"`
	DebugJavascriptAutoAttachFilter                string   `json:"debug.javascript.autoAttachFilter,omitempty"`
	DebugJavascriptAutoAttachSmartPattern          []string `json:"debug.javascript.autoAttachSmartPattern,omitempty"`
	DebugJavascriptAutomaticallyTunnelRemoteServer bool     `json:"debug.javascript.automaticallyTunnelRemoteServer,omitempty"`
	DebugJavascriptBreakOnConditionalError         bool     `json:"debug.javascript.breakOnConditionalError,omitempty"`
	DebugJavascriptCodelensNpmScripts              string   `json:"debug.javascript.codelens.npmScripts,omitempty"`
	DebugJavascriptDebugByLinkOptions              string   `json:"debug.javascript.debugByLinkOptions,omitempty"`
	DebugJavascriptDefaultRuntimeExecutable        struct {
		PwaNode string `json:"pwa-node,omitempty"`
	} `json:"debug.javascript.defaultRuntimeExecutable,omitempty"`
	DebugJavascriptPickAndAttachOptions   struct{} `json:"debug.javascript.pickAndAttachOptions,omitempty"`
	DebugJavascriptResourceRequestOptions struct{} `json:"debug.javascript.resourceRequestOptions,omitempty"`
	DebugJavascriptTerminalOptions        struct{} `json:"debug.javascript.terminalOptions,omitempty"`
	DebugJavascriptUnmapMissingSources    bool     `json:"debug.javascript.unmapMissingSources,omitempty"`
	NpmAutoDetect                         string   `json:"npm.autoDetect,omitempty"`
	NpmEnableRunFromFolder                bool     `json:"npm.enableRunFromFolder,omitempty"`
	NpmEnableScriptExplorer               bool     `json:"npm.enableScriptExplorer,omitempty"`
	NpmExclude                            string   `json:"npm.exclude,omitempty"`
	NpmFetchOnlinePackageInfo             bool     `json:"npm.fetchOnlinePackageInfo,omitempty"`
	NpmPackageManager                     string   `json:"npm.packageManager,omitempty"`
	NpmRunSilent                          bool     `json:"npm.runSilent,omitempty"`
	NpmScriptExplorerAction               string   `json:"npm.scriptExplorerAction,omitempty"`
	NpmScriptExplorerExclude              []any    `json:"npm.scriptExplorerExclude,omitempty"`
	NpmScriptHover                        bool     `json:"npm.scriptHover,omitempty"`
	ReferencesPreferredLocation           string   `json:"references.preferredLocation,omitempty"`
}

import fs from 'node:fs/promises';
import path from 'node:path';
import crypto from 'node:crypto';
import {FileBlob,PresentationFile} from 'file:///C:/Users/iQuantum/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/node_modules/@oai/artifact-tool/dist/artifact_tool.mjs';
import {slides} from './content.mjs';
const root='D:/02-A/code/bootcamp';
const build=root+'/.codex-artifacts/english-demo';
const source=root+'/POS_AI_First_Hackathon_2026.pptx';
const skill='C:/Users/iQuantum/.codex/plugins/cache/openai-primary-runtime/presentations/26.904.11930/skills/presentations';
process.env.RUNTIME_NODE_MODULES='C:/Users/iQuantum/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/node_modules';
const {finalizePresentation}=await import('file:///'+skill+'/container_tools/artifact_tool_utils.mjs');
const p=await PresentationFile.importPptx(await FileBlob.load(source));
if(p.slides.items.length!==22 || slides.length!==22) throw Error('Expected 22 slides');
const inspected=(await p.inspect({kind:'textbox',include:'id,slide,text,bbox',maxChars:100000})).ndjson.trim().split('\n').map(JSON.parse);
const filename=process.argv[2] ?? 'POS_AI_First_Hackathon_2026_EN_Final.pptx';
if(path.basename(filename)!==filename || !filename.endsWith('.pptx')) throw Error('Use a new PPTX basename');
const output=root+'/presentation-output/'+filename;
await fs.mkdir(path.dirname(output),{recursive:true});
for(let i=0;i<22;i++) {
 const slide=p.slides.items[i], spec=slides[i];
 const records=inspected.filter(r=>r.slide===i+1);
 for(const [index,replacement] of Object.entries(spec.edits)) {
  const shape=slide.shapes.items[Number(index)];
  const old=String(shape.text);
  const record=records.find(r=>r.text===old && Math.abs(r.bbox[0]-shape.position.left)<0.1 && Math.abs(r.bbox[1]-shape.position.top)<0.1);
  if(!record) throw Error(`Unverified text target ${i+1}:${index}: ${old}`);
  p.resolve(record.id).text.replace(old,replacement);
 }
 slide.speakerNotes.textFrame.setText(spec.notes);
 const png=await p.export({slide,format:'png',scale:1});
 await fs.writeFile(build+`/english-${i+1}.png`,new Uint8Array(await png.arrayBuffer()));
 await fs.writeFile(build+`/english-${i+1}.layout.json`,await (await slide.export({format:'layout'})).text());
}
const draft=build+'/candidate.pptx';
await (await PresentationFile.exportPptx(p)).save(draft);
const result=await finalizePresentation({workspaceDir:root,candidatePath:draft,finalPath:output,
 explicitTotalSlideCount:22,requiredNativeTableOwnerSlides:[],requiredNativeChartOwnerSlides:[],
 fontPolicy:{basis:'reference',families:['Aptos','Consolas'],referencePath:source,referenceSha256:crypto.createHash('sha256').update(await fs.readFile(source)).digest('hex')},
 pythonExecutable:'C:/Users/iQuantum/.cache/codex-runtimes/codex-primary-runtime/dependencies/python/python.exe',
 integrityValidatorPath:skill+'/container_tools/inspect_presentation_package_integrity.py',
 layoutValidatorPath:skill+'/container_tools/inspect_presentation_layout_geometry.py',
 layoutArgs:['--expected-slide-size-emu','12192000,6858000','--validate-bullet-geometry','--validate-heading-fit'],
 verifyArtifactToolImport:true,receiptPath:build+'/'+filename+'.validation.json'});
console.log(JSON.stringify(result));
const final=await PresentationFile.importPptx(await FileBlob.load(output));
for(let i=0;i<22;i++) {
 const png=await final.export({slide:final.slides.items[i],format:'png',scale:1});
 await fs.writeFile(build+`/final-${i+1}.png`,new Uint8Array(await png.arrayBuffer()));
}
await fs.writeFile(build+'/final-inspect.ndjson',(await final.inspect({kind:'slide,textbox,notes',include:'id,slide,text',maxChars:100000})).ndjson);
const markdown=['# POS AI-First MVP','', '## Bootcamp Kiro × Código Facilito, Hackathon 2026','', '**Presenter:** Gabriel Magallon Sanchez','', 'The app defaults to English and supports Spanish through its language switch. Render deployment remains pending. Keep the existing 22-slide design and use the English PowerPoint alongside these notes.',''];
for(let i=0;i<22;i++) {
 const s=final.slides.items[i];
 const body=s.shapes.items.map(x=>String(x.text)).filter(t=>t.trim() && !['POS AI-FIRST · HACKATHON 2026','KIRO × CÓDIGO FACILITO'].includes(t));
 markdown.push(`## Slide ${i+1}: ${slides[i].title}`,'','**On-slide copy:**','',...body.map(t=>'- '+t.replace(/\n/g,' ')),'','**Speaker notes:** '+slides[i].notes,'','---','');
}
await fs.writeFile(root+'/presentacion.md',markdown.join('\n'));
await fs.writeFile(build+'/source-sha256.txt',crypto.createHash('sha256').update(await fs.readFile(source)).digest('hex')+'\n');
console.log('Saved English deck and aligned presentation notes');

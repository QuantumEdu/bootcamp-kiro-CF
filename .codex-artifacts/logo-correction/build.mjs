import fs from 'node:fs/promises';
import path from 'node:path';
import crypto from 'node:crypto';
import {FileBlob,PresentationFile} from 'file:///C:/Users/iQuantum/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/node_modules/@oai/artifact-tool/dist/artifact_tool.mjs';
const root='D:/02-A/code/bootcamp',build=root+'/.codex-artifacts/logo-correction';
const source=root+'/presentation-output/POS_AI_First_Kiro_University_Showcase_EN_Final.pptx';
const asset=build+'/kiro-user-composite.png',bytes=await fs.readFile(asset);
const sha=b=>crypto.createHash('sha256').update(b).digest('hex');
const sourceSha256=sha(await fs.readFile(source));
if(sha(bytes)!=='06dd54e6f192800f3c83d8778727fd4b18b863dc56f7c0c269b8c28c6c7a4aea') throw Error('User asset changed');
const skill='C:/Users/iQuantum/.codex/plugins/cache/openai-primary-runtime/presentations/26.904.11930/skills/presentations';
process.env.RUNTIME_NODE_MODULES='C:/Users/iQuantum/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/node_modules';
const {finalizePresentation}=await import('file:///'+skill+'/container_tools/artifact_tool_utils.mjs');
const p=await PresentationFile.importPptx(await FileBlob.load(source));
const images=(await p.inspect({kind:'image',include:'id,slide,bbox',maxChars:10000})).ndjson.trim().split('\n').map(JSON.parse);
for(let i of [0,21]) {
 const record=images.find(r=>r.slide===i+1 && Math.abs(r.bbox[2]-116)<.1);
 if(!record) throw Error('Expected logo missing');
 const logo=p.resolve(record.id),old={...logo.frame};
 logo.replace({blob:bytes,contentType:'image/png',alt:'User-supplied purple Kiro word and ghost logo',fit:'contain'});
 const height=116*187/302;
 logo.frame={left:old.left,top:old.top+old.height/2-height/2,width:116,height};
 const notes=p.slides.items[i].speakerNotes.text;
 p.slides.items[i].speakerNotes.setText(notes.replace(/Official (?:Kiro )?logo: https:\/\/kiro\.dev\/images\/kiro-wordmark\.png\?h=0ad65a93\./g,'Logo: user-supplied purple Kiro composite image, embedded unchanged.'));
}
const filename=process.argv[2]??'POS_AI_First_Kiro_University_Showcase_EN_Purple_Logo.pptx';
if(path.basename(filename)!==filename || !filename.endsWith('.pptx')) throw Error('Use new PPTX basename');
const output=root+'/presentation-output/'+filename,draft=build+'/candidate.pptx';
await (await PresentationFile.exportPptx(p)).save(draft);
const result=await finalizePresentation({workspaceDir:root,candidatePath:draft,finalPath:output,explicitTotalSlideCount:22,
 requiredNativeTableOwnerSlides:[],requiredNativeChartOwnerSlides:[],
 fontPolicy:{basis:'reference',families:['Aptos','Consolas'],referencePath:source,referenceSha256:sourceSha256},
 pythonExecutable:'C:/Users/iQuantum/.cache/codex-runtimes/codex-primary-runtime/dependencies/python/python.exe',
 integrityValidatorPath:skill+'/container_tools/inspect_presentation_package_integrity.py',layoutValidatorPath:skill+'/container_tools/inspect_presentation_layout_geometry.py',
 layoutArgs:['--expected-slide-size-emu','12192000,6858000','--validate-bullet-geometry','--validate-heading-fit'],verifyArtifactToolImport:true,receiptPath:build+'/'+filename+'.validation.json'});
const final=await PresentationFile.importPptx(await FileBlob.load(output));
for(let i=0;i<22;i++) {
 const b=await final.export({slide:final.slides.items[i],format:'png',scale:1});
 await fs.writeFile(build+`/final-${i+1}.png`,new Uint8Array(await b.arrayBuffer()));
}
const mdPath=root+'/presentacion.md',md=await fs.readFile(mdPath,'utf8');
await fs.writeFile(mdPath,md.replace('Official logo source: [Kiro wordmark](https://kiro.dev/images/kiro-wordmark.png?h=0ad65a93).','Logo source: user-supplied purple Kiro word and ghost composite, embedded unchanged.').replace(/Official (?:Kiro )?logo: https:\/\/kiro\.dev\/images\/kiro-wordmark\.png\?h=0ad65a93\./g,'Logo: user-supplied purple Kiro composite image, embedded unchanged.'));
await fs.writeFile(build+'/provenance.json',JSON.stringify({source,sourceSha256,asset:'kiro-user-composite.png',origin:'User attachment',assetSha256:sha(bytes),originalDimensions:[302,187],output},null,2)+'\n');
console.log(JSON.stringify(result));

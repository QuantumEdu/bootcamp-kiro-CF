import fs from 'node:fs/promises';
import path from 'node:path';
import crypto from 'node:crypto';
import {FileBlob,PresentationFile} from 'file:///C:/Users/iQuantum/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/node_modules/@oai/artifact-tool/dist/artifact_tool.mjs';
import {slides as originalNotes} from '../english-demo/content.mjs';
const root='D:/02-A/code/bootcamp';
const build=root+'/.codex-artifacts/university-showcase';
const source=root+'/presentation-output/POS_AI_First_Hackathon_2026_EN_Final.pptx';
const logo=build+'/kiro-official-wordmark.png';
const url='https://kiro.dev/2026/university/';
const logoURL='https://kiro.dev/images/kiro-wordmark.png?h=0ad65a93';
const skill='C:/Users/iQuantum/.codex/plugins/cache/openai-primary-runtime/presentations/26.904.11930/skills/presentations';
const hash=b=>crypto.createHash('sha256').update(b).digest('hex');
const sourceHash=hash(await fs.readFile(source));
const logoBytes=await fs.readFile(logo);
if(hash(logoBytes)!=='402933850fe5d2dce859ff3fe30dba89466a203bee63c0e902f5f28d43af78c9') throw Error('Official logo hash mismatch');
process.env.RUNTIME_NODE_MODULES='C:/Users/iQuantum/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/node_modules';
const {finalizePresentation}=await import('file:///'+skill+'/container_tools/artifact_tool_utils.mjs');
const p=await PresentationFile.importPptx(await FileBlob.load(source));
if(p.slides.items.length!==22) throw Error('Expected 22 source slides');
const snapshot=(await p.inspect({kind:'textbox',include:'id,slide,text,bbox',maxChars:100000})).ndjson.trim().split('\n').map(JSON.parse);
const specs={
 1:{title:'POS AI-First learning showcase',edits:{3:'SHOWCASE 2026',15:'Kiro University learning showcase'},notes:'I am Gabriel Magallon Sanchez. This is an educational showcase of an existing POS project and documented learning with Kiro. It is not an eligible final-exam submission. The original bootcamp project predates the University challenge. '+`Official program context: ${url}. Official logo: ${logoURL}.`},
 4:{title:'Kiro University learning framework',edits:{3:'03 · UNIVERSITY',4:'Kiro University learning framework',5:'Official course structure. This showcase does not claim lesson completion.',9:'Study',10:'Lessons',13:'Apply',14:'POS project',17:'Document',18:'Evidence',21:'Verify',22:'Behavior',25:'Share',26:'Learning',28:'7',29:'Required lessons',31:'2',32:'Optional bonuses',34:'3',35:'IDE / CLI / Web',37:'5,250',38:'Potential credits'},notes:'The official program includes seven required lessons and two optional bonuses. Some lessons use the IDE, CLI or Web specifically. The maximum possible award is 5,250 credits, not an award this project has earned. The final-exam deadline is October 5, 2026 at 23:59 PDT. No lesson names or completion status are inferred. '+`Source: ${url}.`},
 5:{title:'Educational showcase scope',edits:{3:'04 · SHOWCASE',4:'An existing project, an educational showcase',5:'A learning demonstration with original history intact.',8:'HISTORY',9:'Existing project',10:'Local start',11:'July 21, 2026',12:'History',13:'Unchanged',16:'PROGRAM',17:'New-build rule',18:'First commit',19:'Sept 21 or later',20:'Context',21:'Challenge',24:'SHOWCASE',26:'Learning',27:'Documented',28:'Eligibility',29:'Not claimed',30:'Existing POS case study',31:'Educational use only. No final-exam eligibility or credit award claimed.'},notes:'The local repository begins July 21, 2026. The official exam requires a new project built during the challenge and a first GitHub commit on or after September 21. This showcase does not satisfy or claim that new-project requirement. Creating a repository or rewriting dates would not turn earlier work into a new build. No history changes are part of this presentation update. '+`Source: ${url}.`},
 7:{title:'Documented Kiro workflow',edits:{3:'06 · LEARNING',4:'Documented Kiro workflow',26:'Project artifacts support the learning story'},notes:originalNotes[6].notes+' These are documented project practices, not a mapping that proves all University lessons complete. '+`Program source: ${url}.`},
 11:{title:'Kiro surfaces and available evidence',edits:{3:'10 · SURFACES',4:'Kiro surfaces and available evidence',6:'I',7:'IDE',8:'Project artifacts',11:'C',12:'CLI',13:'Separate evidence',16:'W',17:'WEB',18:'Separate evidence',20:'SURFACE USE NEEDS ITS OWN EVIDENCE'},notes:'The University page distinguishes IDE, CLI and Web lessons. This project contains IDE-oriented artifacts. CLI or Web completion needs actual evidence, not a diagram. No cross-device synchronization or completed University lesson is claimed here. '+`Source: ${url}.`},
 17:{title:'Lessons from the project',edits:{3:'16 · REFLECTION',4:'Lessons from an existing project'},notes:originalNotes[16].notes+' These are project reflections, not official University syllabus titles.'},
 22:{title:'Thank you',edits:{13:'Educational showcase 2026'},notes:'Thank you. This educational showcase preserves the existing project history and makes no final-exam eligibility or earned-credit claim. '+`Program source: ${url}. Official Kiro logo: ${logoURL}.`}
};
for(let i=0;i<22;i++) {
 const slide=p.slides.items[i];
 const before=await p.export({slide,format:'png',scale:1});
 await fs.writeFile(build+`/source-${i+1}.png`,new Uint8Array(await before.arrayBuffer()));
 const replace=(shape,text)=>{
  const old=String(shape.text),pos=shape.position;
  const record=snapshot.find(r=>r.slide===i+1 && r.text===old && Math.abs(r.bbox[0]-pos.left)<0.1 && Math.abs(r.bbox[1]-pos.top)<0.1);
  if(!record) throw Error('Unverified shape '+(i+1)+': '+old);
  p.resolve(record.id).text.replace(old,text);
 };
 for(const [index,text] of Object.entries(specs[i+1]?.edits??{})) replace(slide.shapes.items[+index],text);
 for(const shape of slide.shapes.items) if(String(shape.text)==='POS AI-FIRST · HACKATHON 2026') replace(shape,'POS AI-FIRST · LEARNING SHOWCASE');
 if(specs[i+1]) slide.speakerNotes.textFrame.setText(specs[i+1].notes);
 if(i===0 || i===21) {
  // Remove the prior editable ghost approximation and text, then embed the official wordmark.
  const indexes=i===0?[11,12,13,14]:[6,7,8,9];
  const ids=indexes.map(n=>slide.shapes.items[n].id);
  for(const id of ids) slide.shapes.deleteById(id);
  slide.images.add({blob:logoBytes,contentType:'image/png',alt:'Official Kiro wordmark',fit:'contain',position:{left:i===0?356:335,top:i===0?547:77,width:116,height:36.21}});
 }
}
const filename=process.argv[2]??'POS_AI_First_Kiro_University_Showcase_EN.pptx';
if(path.basename(filename)!==filename || !filename.endsWith('.pptx')) throw Error('Use a new PPTX basename');
const output=root+'/presentation-output/'+filename;
const draft=build+'/candidate.pptx';
await (await PresentationFile.exportPptx(p)).save(draft);
const result=await finalizePresentation({workspaceDir:root,candidatePath:draft,finalPath:output,
 explicitTotalSlideCount:22,requiredNativeTableOwnerSlides:[],requiredNativeChartOwnerSlides:[],
 fontPolicy:{basis:'reference',families:['Aptos','Consolas'],referencePath:source,referenceSha256:sourceHash},
 pythonExecutable:'C:/Users/iQuantum/.cache/codex-runtimes/codex-primary-runtime/dependencies/python/python.exe',
 integrityValidatorPath:skill+'/container_tools/inspect_presentation_package_integrity.py',layoutValidatorPath:skill+'/container_tools/inspect_presentation_layout_geometry.py',
 layoutArgs:['--expected-slide-size-emu','12192000,6858000','--validate-bullet-geometry','--validate-heading-fit'],
 verifyArtifactToolImport:true,receiptPath:build+'/'+filename+'.validation.json'});
const final=await PresentationFile.importPptx(await FileBlob.load(output));
const md=['# POS AI-First: Kiro University Learning Showcase','','**Presenter:** Gabriel Magallon Sanchez','','An educational showcase of an existing project, not an eligible University final-exam submission. Original project history remains intact. English is the app default, Spanish is available and Render deployment remains pending.','','## Official program context','',`[Kiro University Challenge](${url}) describes seven required lessons and two optional bonuses, with some lessons specific to IDE, CLI or Web. The maximum possible award is 5,250 credits, not credits earned here. The final exam is due October 5, 2026 at 23:59 PDT and requires a new project whose first GitHub commit is September 21 or later. This project has earlier local history. Lesson completion and eligibility are not claimed.`, '',`Official logo source: [Kiro wordmark](${logoURL}). Brand use identifies the learning tool and does not imply endorsement.`,''];
for(let i=0;i<22;i++) {
 const slide=final.slides.items[i];
 const png=await final.export({slide,format:'png',scale:1});
 await fs.writeFile(build+`/final-${i+1}.png`,new Uint8Array(await png.arrayBuffer()));
 const body=slide.shapes.items.map(x=>String(x.text)).filter(t=>t.trim() && !['POS AI-FIRST · LEARNING SHOWCASE','KIRO × CÓDIGO FACILITO'].includes(t));
 md.push(`## Slide ${i+1}: ${specs[i+1]?.title??originalNotes[i].title}`,'','**On-slide copy:**','',...body.map(t=>'- '+t.replace(/\n/g,' ')),'','**Speaker notes:** '+(specs[i+1]?.notes??originalNotes[i].notes),'','---','');
}
await fs.writeFile(root+'/presentacion.md',md.join('\n'));
await fs.writeFile(build+'/provenance.json',JSON.stringify({source,sourceHash,logoURL,logoSha256:hash(logoBytes),programURL:url,retrieved:'2026-09-27',output},null,2));
console.log(JSON.stringify(result));

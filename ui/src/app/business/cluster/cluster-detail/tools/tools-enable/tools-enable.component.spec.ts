import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ToolsEnableComponent } from './tools-enable.component';

describe('ToolsEnableComponent', () => {
  let component: ToolsEnableComponent;
  let fixture: ComponentFixture<ToolsEnableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ToolsEnableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ToolsEnableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

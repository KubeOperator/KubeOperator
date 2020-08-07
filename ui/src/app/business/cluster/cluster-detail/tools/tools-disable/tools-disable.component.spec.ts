import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ToolsDisableComponent } from './tools-disable.component';

describe('ToolsDisableComponent', () => {
  let component: ToolsDisableComponent;
  let fixture: ComponentFixture<ToolsDisableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ToolsDisableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ToolsDisableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

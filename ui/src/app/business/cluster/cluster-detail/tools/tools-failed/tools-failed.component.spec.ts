import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ToolsFailedComponent } from './tools-failed.component';

describe('ToolsFailedComponent', () => {
  let component: ToolsFailedComponent;
  let fixture: ComponentFixture<ToolsFailedComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ToolsFailedComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ToolsFailedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ToolsUpgradeComponent } from './tools-upgrade.component';

describe('ToolsUpgradeComponent', () => {
  let component: ToolsUpgradeComponent;
  let fixture: ComponentFixture<ToolsUpgradeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ToolsUpgradeComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ToolsUpgradeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
